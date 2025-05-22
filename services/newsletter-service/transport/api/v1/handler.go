package v1

import (
	"encoding/json"
	"net/http"

	"your_project_path/service/newsletter"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type NewsletterHandler struct {
	service newsletter.Service
}

func NewNewsletterHandler(service newsletter.Service) *NewsletterHandler {
	return &NewsletterHandler{service: service}
}

func (h *NewsletterHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/newsletters", h.createNewsletter)
	r.Get("/newsletters", h.listNewsletters)
	r.Put("/newsletters/{id}", h.updateNewsletter)

	return r
}

func (h *NewsletterHandler) createNewsletter(w http.ResponseWriter, r *http.Request) {
	var input newsletter.CreateNewsletterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	n, err := h.service.CreateNewsletter(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

func (h *NewsletterHandler) listNewsletters(w http.ResponseWriter, r *http.Request) {
	ns, err := h.service.ListNewsletters(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch newsletters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ns)
}

func (h *NewsletterHandler) updateNewsletter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input newsletter.UpdateNewsletterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	n, err := h.service.UpdateNewsletter(r.Context(), id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}
