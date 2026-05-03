package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	input := string(data)

	// ★ここ重要：URLエンコード
	u := "http://localhost:8787/push-data?data=" + url.QueryEscape(input)

	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	// stdinの内容をそのままstdoutへ
	_, err = os.Stdout.Write(data)
	if err != nil {
		panic(err)
	}
}