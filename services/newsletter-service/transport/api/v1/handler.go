package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"newsletter-service/service/newsletter"

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
	r.Use(middleware.Recoverer) // Recover from panics

	r.Post("/newsletters", h.createNewsletter)
	r.Get("/newsletters", h.listNewsletters)
	r.Put("/newsletters/{id}", h.updateNewsletter)

	return r
}

func (h *NewsletterHandler) createNewsletter(w http.ResponseWriter, r *http.Request) {
	var input newsletter.CreateNewsletterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
		return
	}

	n, err := h.service.CreateNewsletter(r.Context(), input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error": "failed to create newsletter: %v"}`, err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(n); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *NewsletterHandler) listNewsletters(w http.ResponseWriter, r *http.Request) {
	ns, err := h.service.ListNewsletters(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "failed to fetch newsletters"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ns); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *NewsletterHandler) updateNewsletter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var input newsletter.UpdateNewsletterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
		return
	}

	n, err := h.service.UpdateNewsletter(r.Context(), id, input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error": "failed to update newsletter: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(n); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}
