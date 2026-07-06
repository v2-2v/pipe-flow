package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Data struct {
	Command string `json:"command"`
	Data    string `json:"data"`
}

func main() {
	url := "http://localhost:8787/show"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var list []Data
	if err := json.Unmarshal(body, &list); err == nil {

		printList(list)
		return
	}

	var single Data
	if err := json.Unmarshal(body, &single); err == nil {

		printList([]Data{single})
		return
	}

	fmt.Println("JSON parse failed")
	fmt.Println(string(body))
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

func inferWithLMStudio(command string) {
	env := loadEnv(".env")
	endpoint := env["LMSTUDIO_URL"]
	model := env["LMSTUDIO_MODEL"]
	apiKey := env["API_KEY"]

	if endpoint == "" || model == "" {
		fmt.Println("\nSkipped LLM inference: missing .env settings")
		return
	}

	prompt := fmt.Sprintf("Please explain the following shell command in Japanese in one short sentence:\n%s", command)
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
		fmt.Printf("\nFailed to build LLM request: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("\nFailed to create LLM request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	// ローディング画面開始
	fmt.Print("\n推論中")
	done := make(chan struct{})
	go showLoadingSpinner(done)

	client := &http.Client{}
	resp, err := client.Do(req)
	close(done)
	fmt.Print("\r          \r") // ローディング表示をクリア

	if err != nil {
		fmt.Printf("LLM request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read LLM response: %v\n", err)
		return
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		fmt.Printf("Failed to parse LLM response: %v\n", err)
		fmt.Printf("Response: %s\n", string(responseBody))
		return
	}

	if len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		fmt.Println("LLM response was empty")
		return
	}

	fmt.Printf("説明：%s\n", strings.TrimSpace(result.Choices[0].Message.Content))
}

func printList(data []Data) {
	if len(data) == 0 {
		fmt.Println("No data received")
		return
	}

	for i, v := range data {
		d := v.Data
		c := v.Command
		c = strings.TrimLeft(c, " ")
		c = strings.TrimRight(c, " ")
		switch d {
		case "init":
			d = "init (not input data yet)"
		case "":
			d = "nil (empty string)"
		}

		if i == len(data)-1 {
			fmt.Printf("[%s]\n↓\n%s", c, d)
		} else {
			fmt.Printf("[%s]\n↓\n%s\n↓\n", c, d)
		}
	}
	command := ""

	for _, v := range data {
		c := v.Command
		c = strings.TrimLeft(c, " ")
		c = strings.TrimRight(c, " ")
		c += " | "
		command += c
	}
	r := []rune(command)
	if len(r) >= 3 {
		command = string(r[:len(r)-3])
	}
	fmt.Printf("\n\nInput command: %s\n", command)
	inferWithLMStudio(command)
}

func showLoadingSpinner(done <-chan struct{}) {
	spinner := []string{"|", "/", "-", "\\"}
	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fmt.Printf("\r推論中 %s", spinner[i%len(spinner)])
			i++
		}
	}
}
