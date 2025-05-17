 package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/abiopereira/sw-criciuma/internal/domain"
)

type MessageHandler struct {
	messageUseCase domain.MessageUseCase
}

func NewMessageHandler(router *mux.Router, useCase domain.MessageUseCase) {
	handler := &MessageHandler{
		messageUseCase: useCase,
	}

	router.HandleFunc("/api/messages", handler.GetAll).Methods("GET")
	router.HandleFunc("/api/messages", handler.Create).Methods("POST")
}

func (h *MessageHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	messages, err := h.messageUseCase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Messages []domain.Message `json:"messages"`
	}{
		Messages: messages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var message domain.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.messageUseCase.Create(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}{
		ID:      message.ID,
		Message: "Mensagem enviada com sucesso",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
} 