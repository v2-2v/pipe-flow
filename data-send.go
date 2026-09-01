package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func normalizeData(raw string) string {
	raw = strings.Trim(raw, "\r\n")
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\r\n")
}

func main() {
	index := -1
	if len(os.Args) > 1 {
		parsed, err := strconv.Atoi(os.Args[1])
		if err == nil {
			index = parsed
		}
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	dataText := normalizeData(string(data))

	payload := struct {
		Data  string `json:"data"`
		Index int    `json:"index"`
	}{
		Data:  dataText,
		Index: index,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		"http://localhost:8787/push-data",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, err = os.Stdout.WriteString(dataText)
	if err != nil {
		panic(err)
	}
}
