package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	body := struct {
		Data string `json:"data"`
	}{
		Data: string(data),
	}

	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		"http://localhost:8787/push-data",
		"application/json",
		bytes.NewBuffer(b),
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