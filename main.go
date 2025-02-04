package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idVars := vars["id"]
	id, err := strconv.Atoi(idVars)
	if err != nil {
		http.Error(w, "Неверный формат id", http.StatusBadRequest)
		return
	}

	var updatedMessage Message

	if err := json.NewDecoder(r.Body).Decode(&updatedMessage); err != nil {
		http.Error(w, "Неверный JSON формат", http.StatusBadRequest)
		return
	}

	if updatedMessage.Task == "" && updatedMessage.IsDone == nil {
		http.Error(w, "JSON должен содержать хотя бы одно поле: 'task' или 'is_done'", http.StatusBadRequest)
		return
	}

	result := DB.Model(&Message{}).Where("id = ?", id).Updates(updatedMessage)

	if result.Error != nil {
		http.Error(w, "Ошибка при обновлении записи", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Запись не найдена", http.StatusNotFound)
		return
	}

	var message Message

	if err := DB.Take(&message, id).Error; err != nil {
		http.Error(w, "Ошибка при получении обновленной записи", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func main() {
	InitDB()

	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{id}", UpdateMessage).Methods("PATCH")
	http.ListenAndServe(":8080", router)
}
