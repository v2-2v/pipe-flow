package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"fmt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("argument is required")
		return
	}

	line := os.Args[1]
	args := strings.Split(line, "|")

	body := struct {
		Commands []string `json:"commands"`
	}{
		Commands: args,
	}

	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8787/push-command",
		bytes.NewBuffer(b),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}