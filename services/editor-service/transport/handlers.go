package transport

import (
	"context"
	"editor-service/service"
	"encoding/json"
	"net/http"
)

type EditorHandler struct {
	service *service.EditorService
}

func NewEditorHandler(service *service.EditorService) *EditorHandler {
	return &EditorHandler{service: service}
}

type signUpRequest struct {
	IDToken   string `json:"id_token"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *EditorHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.SignUp(context.Background(), req.IDToken, req.FirstName, req.LastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Editor created successfully",
	})
}
