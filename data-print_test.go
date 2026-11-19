package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestRenderTableWrapsLongCommandsToTerminalWidth(t *testing.T) {
	m := Model{
		width: 40,
		dataList: []Data{
			{Command: "cat access.log", Data: "init"},
			{Command: `grep "192.168.1.10"`, Data: "init"},
			{Command: "sort", Data: "init"},
			{Command: "uniq -c", Data: "init"},
			{Command: "sort -nr", Data: "init"},
			{Command: "head -10", Data: "init"},
		},
	}

	rendered := m.renderTable()
	for _, line := range strings.Split(rendered, "\n") {
		if width := lipgloss.Width(line); width > m.width {
			t.Errorf("rendered line width %d exceeds terminal width %d: %q", width, m.width, line)
		}
	}

	for _, commandPart := range []string{"access.log", "192.168.1.10", "uniq", "head"} {
		if !strings.Contains(rendered, commandPart) {
			t.Errorf("rendered output is missing command part %q", commandPart)
		}
	}
}

func TestRenderTablePreservesPipelineAndPreviewsEveryOutput(t *testing.T) {
	m := Model{
		width:       80,
		selectedIdx: 1,
		dataList: []Data{
			{Command: "first", Data: "first output"},
			{Command: "second", Data: "line 1\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7"},
			{Command: "third", Data: "third output"},
		},
	}

	rendered := m.renderTable()
	if !strings.Contains(rendered, "line 1") ||
		!strings.Contains(rendered, "line 2") ||
		!strings.Contains(rendered, "line 3") ||
		!strings.Contains(rendered, "line 4") ||
		!strings.Contains(rendered, "line 5") ||
		!strings.Contains(rendered, "... (2 more lines; press m)") {
		t.Errorf("output was not limited to five lines with an omitted-lines indicator:\n%s", rendered)
	}
	if strings.Contains(rendered, "line 6") || strings.Contains(rendered, "line 7") {
		t.Errorf("preview includes lines beyond the five-line limit:\n%s", rendered)
	}
	if !strings.Contains(rendered, "[1] first") ||
		!strings.Contains(rendered, "[2] second") ||
		!strings.Contains(rendered, "[3] third") ||
		!strings.Contains(rendered, "first output") ||
		!strings.Contains(rendered, "third output") {
		t.Errorf("rendered output does not preserve the full pipeline flow:\n%s", rendered)
	}
}

func TestMKeyTogglesFullOutputForSelectedStep(t *testing.T) {
	m := Model{
		selectedIdx: 1,
		dataList: []Data{
			{Command: "first", Data: "first output"},
			{Command: "second", Data: "line 1\nline 2\nline 3\nline 4\nline 5\nline 6"},
		},
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}}
	updated, _ := m.Update(msg)
	expanded := updated.(Model)
	if !expanded.showFullOutput || !strings.Contains(expanded.renderTable(), "line 6") {
		t.Fatal("m should show the full output of the selected step")
	}

	updated, _ = expanded.Update(msg)
	collapsed := updated.(Model)
	if collapsed.showFullOutput || strings.Contains(collapsed.renderTable(), "line 6") {
		t.Fatal("m should return the selected step to its five-line preview")
	}
}

func TestRenderGlobalInferenceFitsNarrowTerminal(t *testing.T) {
	m := Model{
		width:           32,
		globalInference: "このコマンドはログを抽出して件数順に表示します。",
		dataList: []Data{
			{Command: "cat access.log"},
			{Command: `grep "192.168.1.10"`},
			{Command: "sort"},
			{Command: "uniq -c"},
			{Command: "sort -nr"},
			{Command: "head -10"},
		},
	}

	rendered := m.renderGlobalInference()
	for _, line := range strings.Split(rendered, "\n") {
		if width := lipgloss.Width(line); width > m.width {
			t.Errorf("rendered line width %d exceeds terminal width %d: %q", width, m.width, line)
		}
	}
}

func TestPipelineCommandJoinsCommandsWithoutTruncation(t *testing.T) {
	data := []Data{
		{Command: "cat access.log"},
		{Command: `grep "192.168.1.10"`},
		{Command: "sort"},
		{Command: "uniq -c"},
		{Command: "sort -nr"},
		{Command: "head -10"},
	}
	want := `cat access.log | grep "192.168.1.10" | sort | uniq -c | sort -nr | head -10`

	if got := pipelineCommand(data); got != want {
		t.Errorf("pipelineCommand() = %q, want %q", got, want)
	}
}
