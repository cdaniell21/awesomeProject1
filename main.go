package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message

	if err := DB.Find(&messages).Error; err != nil {
		http.Error(w, "Ошибка при получении записей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(messages)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message Message

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Неверный JSON формат", http.StatusBadRequest)
		return
	}

	if message.Task == "" || message.IsDone == nil {
		http.Error(w, "JSON должен содержать поля 'task' и 'is_done'", http.StatusBadRequest)
		return
	}

	if err := DB.Create(&message).Error; err != nil {
		http.Error(w, "Ошибка при создании записи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func main() {
	InitDB()

	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	http.ListenAndServe(":8080", router)
}
