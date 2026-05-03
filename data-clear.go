package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("http://localhost:8787/clear")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Data cleared successfully")
}