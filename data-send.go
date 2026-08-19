package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	index := -1
	if len(os.Args) > 1 {
		parsed, err := strconv.Atoi(os.Args[1])
		if err == nil {
			index = parsed
		}
	}

	data, err := io.ReadAll(os.Stdin)
	data = bytes.TrimRight(data, "\r\n")
	if err != nil {
		panic(err)
	}

	payload := struct {
		Data  string `json:"data"`
		Index int    `json:"index"`
	}{
		Data:  string(data),
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

	_, err = os.Stdout.Write(data)
	if err != nil {
		panic(err)
	}
}
