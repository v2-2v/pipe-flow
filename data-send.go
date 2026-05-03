package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	u := "http://localhost:8787/push-data"

	// ★ここでデータをPOSTに入れる
	resp, err := http.Post(
		u,
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	// stdoutへそのまま流す
	_, err = os.Stdout.Write(data)
	if err != nil {
		panic(err)
	}
}