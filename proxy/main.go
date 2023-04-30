package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
)

var FormatToHost map[string]string


type Responce struct {
	Result string `json:"result"`
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
    host, ok := FormatToHost[format]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
    

	res, err := http.Get(host + "/get_result")
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    resp := Responce{}
	json.Unmarshal(body, &resp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("incorrect num of args provided: 2 is required")
	}
	port := args[1]
	
	FormatToHost = make(map[string]string)
    FormatToHost["json"] = "http://0.0.0.0:2001"
	FormatToHost["xml"] = "http://0.0.0.0:2002"
	FormatToHost["yaml"] = "http://0.0.0.0:2003"
	
	http.HandleFunc("/get_result", GetResult)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("server stopped with error: %s", err)
	}
}
