package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

	// 文字列配列として受け取る
	var data []string
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	if data == nil {
		fmt.Println("No data received")
		return
	}
	// forで回す
	for i, v := range data {
		fmt.Println(i, v)
	}
}