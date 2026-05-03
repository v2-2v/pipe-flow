package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Post("http://localhost:8787/clear", "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Data cleared successfully")
}