package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("argument is required")
		return
	}

	data := os.Args[1]
	fmt.Printf("Sending data: %s\n", data)
	resp, err := http.Get("http://localhost:8787/push-data?data=" + data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}