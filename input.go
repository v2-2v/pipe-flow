package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("argument is required")
		return
	}

	// 2個目以降すべて取得
	line := os.Args[1]
	args := strings.Split(line, "|")
	commands := ""
	for _, command := range args {
		// URLエンコード推奨
		commands += "&command=" + url.QueryEscape(command)
	}
	u := "http://localhost:8787/push-command?" + commands
	resp, err := http.Post(u, "application/json", nil)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}