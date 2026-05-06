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

func pushcommandHandler(w http.ResponseWriter, r *http.Request) { //POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Commands []string `json:"commands"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(req.Commands) == 0 {
		http.Error(w, "Missing commands", http.StatusBadRequest)
		return
	}
	for _, command := range req.Commands {
		datas = append(datas, Data{
			Command: command,
			Data:     "init",
		})
	}
	fmt.Printf("new datas: %+v\n", datas)
	json.NewEncoder(w).Encode(Response{
		Message: "push received",
	})
}

func pushdataHandler(w http.ResponseWriter, r *http.Request) { //POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received data: %s\n", req.Data)
	for i := range datas {
		if datas[i].Data == "init" {
			datas[i].Data = req.Data

			fmt.Printf("new datas: %+v\n", datas)

			json.NewEncoder(w).Encode(Response{
				Message: "push received",
			})
			return
		}
	}
	http.Error(w, "No command to associate with data", http.StatusBadRequest)
}

func clearHandler(w http.ResponseWriter, r *http.Request) { //POST
	if r.Method != http.MethodPost {
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

func showHandler(w http.ResponseWriter, r *http.Request) { //GET
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