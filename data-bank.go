package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
)

type Data struct {
	Input  string `json:"input"`
	Command string `json:"command"`
	Output string `json:"output"`
}

var datas []Data

type Response struct {
	Message string `json:"message"`
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	command := r.URL.Query().Get("command")
	input := r.URL.Query().Get("input")
	output := r.URL.Query().Get("output")
	if command == "" {
		http.Error(w, "missing command, input, or output parameter", http.StatusBadRequest)
		return
	}
	datas = append(datas, Data{
		Input:  input,
		Command: command,
		Output: output,
	})
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
	datas = []Data{}
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

	log.Println("Server started at :8787")
	log.Fatal(http.ListenAndServe(":8787", nil))
}