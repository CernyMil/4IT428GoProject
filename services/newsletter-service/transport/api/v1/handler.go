package v1

import (
	"encoding/json"
	"net/http"

	"newsletter-service/repository"
	"newsletter-service/service/model"
	"newsletter-service/transport/middleware"

	"github.com/go-chi/chi/v5"
)

// Handler handles newsletter-related HTTP requests.
type Handler struct {
	*chi.Mux

	authenticator middleware.FirebaseAuthenticator
	service       model.Service
}

type NewsletterHandler struct {
	repo repository.Repository
}

func NewNewsletterHandler(repo repository.Repository) *NewsletterHandler {
	return &NewsletterHandler{repo: repo}
}

func (h *NewsletterHandler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/newsletters", h.createNewsletter)
	r.Get("/newsletters", h.listNewsletters)
	return r
}

func (h *NewsletterHandler) createNewsletter(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a newsletter
}

func (h *NewsletterHandler) listNewsletters(w http.ResponseWriter, r *http.Request) {
	// Implementation for listing newsletters
}

// NewHandler creates a new instance of Handler.
func NewHandler(
	authenticator middleware.FirebaseAuthenticator,
	service model.Service,
) *Handler {
	h := &Handler{
		authenticator: authenticator,
		service:       service,
	}
	h.initRouter()
	return h
}

// initRouter sets up the routes and middleware for the handler.
func (h *Handler) initRouter() {
	r := chi.NewRouter()

	// Setup middleware
	authenticate := h.authenticator.Authenticate

	// Newsletter routes
	r.Route("/newsletters", func(r chi.Router) {
		r.Post("/", h.CreateNewsletter)
		r.Get("/", h.ListNewsletters)
		r.With(authenticate).Put("/{id}", h.UpdateNewsletter)
		r.With(authenticate).Delete("/{id}", h.DeleteNewsletter)

		// Post routes
		r.Route("/{id}/posts", func(r chi.Router) {
			r.Post("/", h.CreatePost)
			r.Get("/", h.ListPosts)
			r.With(authenticate).Put("/{postID}", h.UpdatePost)
			r.With(authenticate).Delete("/{postID}", h.DeletePost)
		})
	})

	h.Mux = r
}

// CreateNewsletter handles the creation of a new newsletter.
func (h *Handler) CreateNewsletter(w http.ResponseWriter, r *http.Request) {
	var input model.CreateNewsletterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
		return
	}

	n, err := h.service.CreateNewsletter(r.Context(), input)
	if err != nil {
		http.Error(w, `{"error": "failed to create newsletter"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

// ListNewsletters handles listing all newsletters.
func (h *Handler) ListNewsletters(w http.ResponseWriter, r *http.Request) {
	newsletters, err := h.service.ListNewsletters(r.Context())
	if err != nil {
		http.Error(w, `{"error": "failed to fetch newsletters"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsletters)
}

// UpdateNewsletter handles updating an existing newsletter.
func (h *Handler) UpdateNewsletter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var input model.UpdateNewsletterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
		return
	}

	n, err := h.service.UpdateNewsletter(r.Context(), id, input)
	if err != nil {
		http.Error(w, `{"error": "failed to update newsletter"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}

// DeleteNewsletter handles deleting a newsletter.
func (h *Handler) DeleteNewsletter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteNewsletter(r.Context(), id); err != nil {
		http.Error(w, `{"error": "failed to delete newsletter"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreatePost handles the creation of a new post.
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "id")
	if newsletterID == "" {
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var input model.CreatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
		return
	}

	p, err := h.service.CreatePost(r.Context(), newsletterID, input)
	if err != nil {
		http.Error(w, `{"error": "failed to create post"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// ListPosts handles listing all posts for a specific newsletter.
func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "id")
	if newsletterID == "" {
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	posts, err := h.service.ListPosts(r.Context(), newsletterID)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch posts"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// UpdatePost handles updating an existing post.
func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "id")
	postID := chi.URLParam(r, "postID")
	if newsletterID == "" || postID == "" {
		http.Error(w, `{"error": "missing newsletter or post ID"}`, http.StatusBadRequest)
		return
	}

	var input model.UpdatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
		return
	}

	p, err := h.service.UpdatePost(r.Context(), newsletterID, postID, input)
	if err != nil {
		http.Error(w, `{"error": "failed to update post"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// DeletePost handles deleting a post.
func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	newsletterID := chi.URLParam(r, "id")
	postID := chi.URLParam(r, "postID")
	if newsletterID == "" || postID == "" {
		http.Error(w, `{"error": "missing newsletter or post ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePost(r.Context(), newsletterID, postID); err != nil {
		http.Error(w, `{"error": "failed to delete post"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
