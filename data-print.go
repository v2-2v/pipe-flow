package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Data struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

type Model struct {
	dataList             []Data
	selectedIdx          int
	inference            string
	inferenceSelectedIdx int
	loading              bool
	err                  string
	width                int
	height               int
	showInference        bool
	globalInference      string
	globalLoading        bool
	spinnerIndex         int
	inferenceCtx         context.Context
	cancelInference      context.CancelFunc
}

type inferenceCompleteMsg struct {
	inference string
	isGlobal  bool
}

type inferenceErrorMsg struct {
	err      error
	isGlobal bool
}

type startInferenceMsg struct {
	ctx      context.Context
	cancel   context.CancelFunc
	isGlobal bool
}

type tickMsg struct{}

func (t tickMsg) After() time.Duration {
	return 100 * time.Millisecond
}

func (t tickMsg) String() string {
	return "tick"
}

func normalizeData(raw string) string {
	raw = strings.Trim(raw, "\r\n")
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\r\n")
}

func indentBlock(text string) string {
	text = normalizeData(text)
	if text == "" {
		return ""
	}
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = "  " + line
	}
	return strings.Join(lines, "\n")
}

func main() {
	m := Model{
		selectedIdx:   0,
		globalLoading: true,
	}

	// データ取得
	if err := m.fetchData(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch data: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m *Model) fetchData() error {
	url := "http://localhost:8787/show"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var list []Data
	if err := json.Unmarshal(body, &list); err == nil {
		m.dataList = list
		return nil
	}

	var single Data
	if err := json.Unmarshal(body, &single); err == nil {
		m.dataList = []Data{single}
		return nil
	}

	return fmt.Errorf("JSON parse failed")
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.buildGlobalInferenceCmd(),
		m.tick(),
	)
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.selectedIdx > 0 {
				if m.cancelInference != nil {
					m.cancelInference()
				}
				m.selectedIdx--
				m.inference = ""
				m.showInference = false
				m.loading = false
			}
		case "down":
			if m.selectedIdx < len(m.dataList)-1 {
				if m.cancelInference != nil {
					m.cancelInference()
				}
				m.selectedIdx++
				m.inference = ""
				m.showInference = false
				m.loading = false
			}
		case "l":
			if !m.loading && len(m.dataList) > 0 {
				m.loading = true
				m.inference = ""
				m.inferenceSelectedIdx = m.selectedIdx
				cmd := m.buildInferenceCmd()
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case inferenceCompleteMsg:
		if msg.isGlobal {
			m.globalInference = msg.inference
			m.globalLoading = false
		} else {
			// selectedIdxが変わっていない場合だけ結果を設定
			if m.inferenceSelectedIdx == m.selectedIdx {
				m.inference = msg.inference
				m.showInference = true
			}
			m.loading = false
		}
	case inferenceErrorMsg:
		if msg.isGlobal {
			m.globalLoading = false
		} else {
			m.err = msg.err.Error()
			m.loading = false
		}
	case tickMsg:
		if m.globalLoading || m.loading {
			m.spinnerIndex++
		}
		return m, m.tick()
	}
	return m, nil
}

func (m Model) View() string {
	if len(m.dataList) == 0 {
		return "No data received"
	}

	var sb strings.Builder

	// ヘッダー
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("4")).
		Padding(0, 1)

	sb.WriteString(headerStyle.Render("PIPE-FLOW DATA VIEWER"))
	sb.WriteString("\n\n")

	// グローバル推論結果
	if m.globalLoading || m.globalInference != "" {
		sb.WriteString(m.renderGlobalInference())
		sb.WriteString("\n\n")
	}

	// テーブル
	sb.WriteString(m.renderTable())
	sb.WriteString("\n")

	// 推論結果表示
	if m.showInference || m.loading {
		sb.WriteString(m.renderInference())
		sb.WriteString("\n")
	}

	// フッター
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))
	sb.WriteString(footerStyle.Render(
		"↑/↓: Select | l: Explain | q: Quit",
	))

	if m.err != "" {
		errStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1"))
		sb.WriteString("\n")
		sb.WriteString(errStyle.Render("Error: " + m.err))
	}

	return sb.String()
}

func (m Model) renderTable() string {
	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("4")).
		Foreground(lipgloss.Color("15"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	var sb strings.Builder

	for i, item := range m.dataList {
		cmd := strings.TrimSpace(item.Command)
		data := normalizeData(item.Data)
		if data == "init" {
			data = "init (not input data yet)"
		} else if data == "" {
			data = "nil (empty string)"
		} else {
			data = indentBlock(data)
		}

		style := normalStyle
		if i == m.selectedIdx {
			style = selectedStyle
		}

		prefix := " "
		if i == m.selectedIdx {
			prefix = ">"
		}

		sb.WriteString(style.Render(fmt.Sprintf("%s [%s]", prefix, cmd)))
		sb.WriteString("\n")
		sb.WriteString(style.Render(fmt.Sprintf("  ↓")))
		sb.WriteString("\n")
		sb.WriteString(style.Render(data))
		sb.WriteString("\n")

		if i < len(m.dataList)-1 {
			sb.WriteString(normalStyle.Render("  ↓"))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m Model) renderGlobalInference() string {
	width := m.width - 4
	if width < 40 {
		width = 40
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("2")).
		Padding(1).
		Width(width)

	if m.globalLoading {
		spinner := []string{"|", "/", "-", "\\"}
		content := fmt.Sprintf("Loading... %s", spinner[m.spinnerIndex%len(spinner)])
		return boxStyle.Render(content)
	}

	if m.globalInference != "" {
		// パイプラインコマンド構築
		command := ""
		for _, v := range m.dataList {
			c := strings.TrimSpace(v.Command)
			command += c + " | "
		}
		r := []rune(command)
		if len(r) >= 3 {
			command = string(r[:len(r)-3])
		}

		content := fmt.Sprintf("Global Command Explanation\n\n%s\n\n%s", command, m.globalInference)
		return boxStyle.Render(content)
	}

	return ""
}

func (m Model) renderInference() string {
	width := m.width - 4
	if width < 40 {
		width = 40
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("4")).
		Padding(1).
		Width(width)

	if m.loading {
		spinner := []string{"|", "/", "-", "\\"}
		content := fmt.Sprintf("Loading... %s", spinner[m.spinnerIndex%len(spinner)])
		return boxStyle.Render(content)
	}

	if m.inference != "" {
		command := strings.TrimSpace(m.dataList[m.inferenceSelectedIdx].Command)
		content := fmt.Sprintf("Command Explanation\n\n%s\n\n%s", command, m.inference)
		return boxStyle.Render(content)
	}

	return ""
}

func (m Model) buildGlobalInferenceCmd() tea.Cmd {
	return func() tea.Msg {
		command := ""
		for _, v := range m.dataList {
			c := strings.TrimSpace(v.Command)
			command += c + " | "
		}
		r := []rune(command)
		if len(r) >= 3 {
			command = string(r[:len(r)-3])
		}

		var dataParts []string
		for i, v := range m.dataList {
			data := strings.TrimSpace(v.Data)
			if data == "" || data == "init" {
				continue
			}
			dataParts = append(dataParts, fmt.Sprintf("Step %d command: %s\nData: %s", i+1, strings.TrimSpace(v.Command), data))
		}
		dataText := strings.Join(dataParts, "\n\n")
		if dataText == "" {
			dataText = "none"
		}

		inference, err := m.runInferenceWithPrompt(command, dataText, buildGlobalInferencePrompt)
		if err != nil {
			return inferenceErrorMsg{err, true}
		}

		return inferenceCompleteMsg{inference, true}
	}
}

func (m Model) buildInferenceCmd() tea.Cmd {
	return func() tea.Msg {
		if m.selectedIdx < 0 || m.selectedIdx >= len(m.dataList) {
			return inferenceErrorMsg{fmt.Errorf("no command selected"), false}
		}

		command := strings.TrimSpace(m.dataList[m.selectedIdx].Command)
		data := strings.TrimSpace(m.dataList[m.selectedIdx].Data)

		inference, err := m.runInference(command, data)
		if err != nil {
			return inferenceErrorMsg{err, false}
		}

		return inferenceCompleteMsg{inference, false}
	}
}

func buildInferencePrompt(command string, data string) string {
	command = strings.TrimSpace(command)
	data = strings.TrimSpace(data)
	if data == "" {
		data = "none"
	}

	return fmt.Sprintf(
		"Explain this shell command in Japanese as a concise learning guide. Use the following sections:\n1. 何をするコマンドか\n2. 構文\n3. 引数・オプションの説明\n4. 使い方（実行例）\n5. 入出力データの読み方\n\nOutput must be plain text only. Do not use Markdown or any Markdown-like formatting: no headings with #, no bullet markers such as -, *, or +, no numbered-list formatting, no backticks, no code blocks, no bold or italic markers, and no tables. Write section names as ordinary plain text lines. Keep the explanation clear and concise. Include only sections that are relevant to this command, and do not invent options or behavior. Show command examples as plain text lines without backticks or other decoration.\n\nCritical rules:\n- Keep every literal string from the command and data exactly as it appears. Do not translate, localize, paraphrase, or change casing of any file name, command, flag, word, or value.\n- Never replace 'orange' with 'オレンジ', 'apple' with 'りんご', or any other translated form. Keep values like 'orange +3' and 'orange -1' exactly unchanged.\n- Only the surrounding explanation may be in Japanese; the actual command/data text must remain verbatim.\n\nCommand:\n%s\n\nInput/output data:\n%s",
		command,
		data,
	)
}

func buildGlobalInferencePrompt(command string, data string) string {
	command = strings.TrimSpace(command)
	data = strings.TrimSpace(data)
	if data == "" {
		data = "none"
	}

	return fmt.Sprintf(
		"Please explain the following shell command in Japanese in one short sentence while considering the input/output data.\n\nImportant: keep the original file contents, command arguments, and data values exactly as they are. Do not translate or paraphrase the actual text in the data. Only explain what the command does in Japanese.\n\nCommand:\n%s\n\nInput/output data:\n%s",
		command,
		data,
	)
}

func (m Model) runInference(command string, data string) (string, error) {
	return m.runInferenceWithPrompt(command, data, buildInferencePrompt)
}

func (m Model) runInferenceWithPrompt(command string, data string, promptBuilder func(string, string) string) (string, error) {
	env := loadEnv(".env")
	endpoint := env["LMSTUDIO_URL"]
	model := env["LMSTUDIO_MODEL"]
	apiKey := env["API_KEY"]

	if endpoint == "" || model == "" {
		return "LLM設定がありません(.env参照)", nil
	}

	prompt := promptBuilder(command, data)
	payload := map[string]any{
		"model": model,
		"messages": []map[string]string{{
			"role":    "user",
			"content": prompt,
		}},
		"temperature": 0.2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		return "LLMレスポンスが空です", nil
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

func loadEnv(path string) map[string]string {
	env := map[string]string{}
	content, err := os.ReadFile(path)
	if err != nil {
		return env
	}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		env[key] = value
	}

	return env
}
