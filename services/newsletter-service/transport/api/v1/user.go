package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"newsletter-management-api/repository"
	dbmodel "newsletter-management-api/repository/sql/model"
	model "newsletter-management-api/repository/sql/model"
	util "newsletter-management-api/transport/api/util"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func getEmailFromURL(r *http.Request) string {
    email := chi.URLParam(r, "email")
    return email
}

// Handler handles user-related HTTP requests.
type Handler struct {
    repo    *repository.UserRepository
    service *repository.UserService
}

func NewHandler(repo *repository.UserRepository, service *repository.UserService) *Handler {
    return &Handler{repo: repo, service: service}
}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
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
    if err := h.service.CreateUser(r.Context(), user); err != nil {
        util.WriteErrResponse(w, http.StatusInternalServerError, err)
        return
    }
    util.WriteResponse(w, http.StatusCreated, user)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    user, err := h.service.GetUser(r.Context(), getEmailFromURL(r))
    if err != nil {
        util.WriteErrResponse(w, http.StatusNotFound, err)
        return
    }
    util.WriteResponse(w, http.StatusOK, user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    users, err := h.service.ListUsers(r.Context())
    if err != nil {
        util.WriteErrResponse(w, http.StatusInternalServerError, err)
        return
    }
    util.WriteResponse(w, http.StatusOK, users)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
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
    email := user.Email
    userUpdated, err := h.service.UpdateUser(r.Context(), email, user)
    if err != nil {
        util.WriteErrResponse(w, http.StatusInternalServerError, err)
        return
    }
    util.WriteResponse(w, http.StatusOK, userUpdated)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if err := h.service.DeleteUser(r.Context(), getEmailFromURL(r)); err != nil {
        util.WriteErrResponse(w, http.StatusNotFound, err)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}