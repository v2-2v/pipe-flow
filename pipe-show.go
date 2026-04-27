package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
)

var datas []string

type Response struct {
	Message string `json:"message"`
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	data := r.URL.Query().Get("data")

	if data == "" {
		http.Error(w, "missing data parameter", http.StatusBadRequest)
		return
	}
	datas = append(datas, data)
	fmt.Printf("new datas: %s\n", datas)
	res := Response{
		Message: "push received",
	}
	json.NewEncoder(w).Encode(res)
}

func clearHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	datas = []string{}
	fmt.Printf("datas cleared\n")
	res := Response{
		Message: "datas cleared",
	}
	json.NewEncoder(w).Encode(res)
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(datas)
}

func main() {
	http.HandleFunc("/push", pushHandler)
	http.HandleFunc("/clear", clearHandler)
	http.HandleFunc("/show", showHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}