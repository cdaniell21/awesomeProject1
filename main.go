package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var task string

type requestBody struct {
	Message string `json:"message"`
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %s", task)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	body := requestBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "JSON должен содержать поле 'message'", http.StatusBadRequest)
		return
	}
	task = body.Message
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/hello", GetHandler).Methods("GET")
	router.HandleFunc("/api/hello", PostHandler).Methods("POST")
	http.ListenAndServe(":8080", router)
}
