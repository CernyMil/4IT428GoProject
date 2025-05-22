package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	"github.com/jackc/pgx/v5/service/model"
	"github.com/jackc/pgx/v5/transport/util"
)

var validate = validator.New()

func getEmailFromURL(r *http.Request) string {
	email := chi.URLParam(r, "email")
	return email
}

func ValidateBirthDate(bd model.BirthDate) error {
	dateStr := fmt.Sprintf("%04d-%02d-%02d", bd.Year, bd.Month, bd.Day)
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		return err
	}
	return nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	var user model.User
	if err := json.Unmarshal(b, &user); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := validate.Struct(user); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return

	}

	if err := ValidateBirthDate(user.BirthDate); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.CreateUser(r.Context(), user); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	util.WriteResponse(w, http.StatusCreated, user)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetUser(r.Context(), getEmailFromURL(r))
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	util.WriteResponse(w, http.StatusOK, user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.ListUsers(r.Context())
	util.WriteResponse(w, http.StatusOK, users)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	var user model.User
	if err := json.Unmarshal(b, &user); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	email := user.Email
	userUpdated, err := h.service.UpdateUser(r.Context(), email, user)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	util.WriteResponse(w, http.StatusOK, userUpdated)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteUser(r.Context(), getEmailFromURL(r)); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
}
