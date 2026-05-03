package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Data struct {
	Command string `json:"command"`
	Data string `json:"data"`
}

func main() {
	url := "http://localhost:8787/show"

	resp, err := http.Post(url, "application/json", nil)
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

	for i, v := range data {
		command := v.Command
		d := v.Data

		if d == "init" {
			d = "init (not input data yet)"
		} else if d == "" {
			d = "nil (empty string)"
		}

		if i == len(data)-1 { // 最後
			fmt.Printf("[%s] --> %s", command, d)
		} else {
			fmt.Printf("[%s] --> %s --> ", command, d)
		}
	}
	fmt.Println()
}