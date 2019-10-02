package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var doc string = ""

type RequestBody struct {
	Snippet string `json:"snippet"`
}

func snippetHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if err != nil {
		log.Println("Failed to read request body:", err)
		return
	}

	data := RequestBody{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Failed to parse JSON:", err)
		return
	}

	fmt.Println(data)
	doc += data.Snippet
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Fprintf(w, doc)
}

func main() {
	http.HandleFunc("/snippet/", snippetHandler)
	http.HandleFunc("/state/", stateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
