package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	// 👇まず配列として試す
	var list []Data
	if err := json.Unmarshal(body, &list); err == nil {

		printList(list)
		return
	}

	// 👇単体オブジェクトとして試す
	var single Data
	if err := json.Unmarshal(body, &single); err == nil {

		printList([]Data{single})
		return
	}

	fmt.Println("JSON parse failed")
	fmt.Println(string(body))
}

func printList(data []Data) {
	if len(data) == 0 {
		fmt.Println("No data received")
		return
	}

	for i, v := range data {
		d := v.Data

		switch d {
		case "init":
			d = "init (not input data yet)"
		case "":
			d = "nil (empty string)"
		}

		if i == len(data)-1 {
			fmt.Printf("[%s] --> %s", v.Command, d)
		} else {
			fmt.Printf("[%s] --> %s --> ", v.Command, d)
		}
	}

	fmt.Println()
}