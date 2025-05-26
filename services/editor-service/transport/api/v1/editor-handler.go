package transport

import (
	"editor-service/service"
	"editor-service/transport/api/v1/model"
	"encoding/json"
	"net/http"
)

type EditorHandler struct {
	service *service.EditorService
}

func NewEditorHandler(service *service.EditorService) *EditorHandler {
	return &EditorHandler{service: service}
}

func (h *EditorHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req model.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.SignUp(r.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Editor created successfully",
	})
}

func (h *EditorHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req model.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	email, err := h.service.VerifyIDTokenAndGetEmail(r.Context(), req.IDToken)
	if err != nil {
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}
	editor, err := h.service.GetByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(editor)
}

func (h *EditorHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req model.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := h.service.ChangePassword(r.Context(), req.Email, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed"})
}
