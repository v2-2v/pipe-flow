package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type Data struct {
    Command string `json:"command"`
    Data    string `json:"data"`
}

var datas []Data

func storeDataByIndex(list []Data, index int, value string) []Data {
    if len(list) == 0 {
        return []Data{{Data: value}}
    }

    if index >= 0 && index < len(list) {
        list[index].Data = value
        return list
    }

    for i := range list {
        if list[i].Data == "init" {
            list[i].Data = value
            return list
        }
    }

    list = append(list, Data{Data: value})
    return list
}

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
            Data:    "init",
        })
    }

    fmt.Printf("new datas: %+v\n", datas)
    _ = json.NewEncoder(w).Encode(Response{Message: "push received"})
}

func pushdataHandler(w http.ResponseWriter, r *http.Request) { //POST
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        Data  string `json:"data"`
        Index int    `json:"index"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    fmt.Printf("Received data: %s (index=%d)\n", req.Data, req.Index)
    datas = storeDataByIndex(datas, req.Index, req.Data)
    fmt.Printf("new datas: %+v\n", datas)

    _ = json.NewEncoder(w).Encode(Response{Message: "push received"})
}

func clearHandler(w http.ResponseWriter, r *http.Request) { //POST
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    datas = []Data{}
    fmt.Printf("datas cleared\n")
    _ = json.NewEncoder(w).Encode(Response{Message: "datas cleared"})
}

func showHandler(w http.ResponseWriter, r *http.Request) { //GET
    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    _ = json.NewEncoder(w).Encode(datas)
}

func main() {
    http.HandleFunc("/push-command", pushcommandHandler)
    http.HandleFunc("/push-data", pushdataHandler)
    http.HandleFunc("/clear", clearHandler)
    http.HandleFunc("/show", showHandler)

    log.Println("Server started at :8787")
    log.Fatal(http.ListenAndServe(":8787", nil))
}
