package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Data struct {
	Input  string `json:"input"`
	Command string `json:"command"`
	Output string `json:"output"`
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

	// Data構造体のスライスとして受け取る
	var data []Data
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
		input := v.Input
		output := v.Output
		if input == "" {input = "nil"}
		if output == "" {output = "nil"}
		fmt.Println(i+1,"  input:", input, "  command:", v.Command, "  output:", output)
	}
}