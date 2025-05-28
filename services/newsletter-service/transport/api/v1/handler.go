package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	newsletter "newsletter-management-api/service/model"

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

	// Newsletter routes
	r.Post("/newsletters", h.createNewsletter)
	r.Get("/newsletters", h.listNewsletters)
	r.Put("/newsletters/{id}", h.updateNewsletter)

	// Post routes
	r.Post("/newsletters/{id}/posts", h.createPost)
	r.Get("/newsletters/{id}/posts", h.listPosts)
	r.Put("/newsletters/{id}/posts/{postID}", h.updatePost)
	r.Delete("/newsletters/{id}/posts/{postID}", h.deletePost)

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
func (h *NewsletterHandler) createPost(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "id")
	if newsletterID == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var input newsletter.CreatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
		return
	}

	p, err := h.service.CreatePost(r.Context(), newsletterID, input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error": "failed to create post: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *NewsletterHandler) listPosts(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "id")
	if newsletterID == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	posts, err := h.service.ListPosts(r.Context(), newsletterID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "failed to fetch posts"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *NewsletterHandler) updatePost(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "id")
	postID := chi.URLParam(r, "postID")
	if newsletterID == "" || postID == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "missing newsletter or post ID"}`, http.StatusBadRequest)
		return
	}

	var input newsletter.UpdatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
		return
	}

	p, err := h.service.UpdatePost(r.Context(), newsletterID, postID, input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error": "failed to update post: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}
func (h *NewsletterHandler) deletePost(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "id")
	postID := chi.URLParam(r, "postID")
	if newsletterID == "" || postID == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "missing newsletter or post ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePost(r.Context(), newsletterID, postID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error": "failed to delete post: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
