package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
)

type Data struct {
	Command string `json:"command"`
	Data string `json:"data"`
}

var datas []Data

type Response struct {
	Message string `json:"message"`
}

func pushcommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	commands := r.URL.Query()["command"] 
	if commands == nil {
		http.Error(w, "Missing commands parameter", http.StatusBadRequest)
		return
	}

	for _, command := range commands {
		datas = append(datas, Data{
			Command: command,
			Data: "init",
		})
	}
	fmt.Printf("new datas: %s\n", datas)
	res := Response{
		Message: "push received",
	}
	json.NewEncoder(w).Encode(res)
}

func pushdataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	data := r.URL.Query().Get("data")

	for i := range datas {
		if datas[i].Data == "init" {
			datas[i].Data = data
			fmt.Printf("new datas: %+v\n", datas)

			res := Response{
				Message: "push received",
			}
			json.NewEncoder(w).Encode(res)
			return
		}
	}

	http.Error(w, "No command to associate with data", http.StatusBadRequest)

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
	http.HandleFunc("/push-command", pushcommandHandler)
	http.HandleFunc("/push-data", pushdataHandler)
	http.HandleFunc("/clear", clearHandler)
	http.HandleFunc("/show", showHandler)

	log.Println("Server started at :8787")
	log.Fatal(http.ListenAndServe(":8787", nil))
}