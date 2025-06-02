package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	id "newsletter-service/pkg/id"
	"newsletter-service/service"
	"newsletter-service/service/model"
	"newsletter-service/transport/middleware"

	"github.com/go-chi/chi/v5"
)

// Handler handles newsletter-related HTTP requests.
type Handler struct {
	*chi.Mux

	authenticator middleware.FirebaseAuthenticator
	service       service.NewsletterService
}

// NewHandler creates a new instance of Handler.
func NewHandler(
	authenticator middleware.FirebaseAuthenticator,
	service service.NewsletterService,
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

	serviceToken := os.Getenv("SERVICE_TOKEN")
	authenticate := h.authenticator.Authenticate

	r.Route("/newsletters", func(r chi.Router) {
		r.With(authenticate).Post("/", h.CreateNewsletter)
		r.With(authenticate).Get("/", h.ListNewsletters)
		r.With(middleware.InternalOnlyMiddleware(serviceToken)).Get("/internal", h.ListNewsletters)
		r.With(authenticate).Put("/{id}", h.UpdateNewsletter)
		r.With(authenticate).Delete("/{id}", h.DeleteNewsletter)

		r.Route("/{id}/posts", func(r chi.Router) {
			r.With(authenticate).Post("/", h.CreatePost)
			r.With(authenticate).Get("/", h.ListPosts)
			r.With(authenticate).Put("/{postID}", h.UpdatePost)
			r.With(authenticate).Delete("/{postID}", h.DeletePost)
			r.With(authenticate).Post("/{postID}/publish", h.PublishPost)
		})
	})

	h.Mux = r
}

func (h *Handler) Routes() http.Handler {
	return h.Mux
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
		fmt.Println("CreateNewsletter error:", err)
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
		fmt.Println("ListNewsletters error:", err)
		http.Error(w, `{"error": "failed to fetch newsletters"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsletters)
}

// UpdateNewsletter handles updating an existing newsletter.
func (h *Handler) UpdateNewsletter(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(idStr); err != nil {
		http.Error(w, `{"error": "invalid newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var input model.UpdateNewsletterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "invalid input"}`, http.StatusBadRequest)
		return
	}

	n, err := h.service.UpdateNewsletter(r.Context(), newsletterID, input)
	if err != nil {
		http.Error(w, `{"error": "failed to update newsletter"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}

// DeleteNewsletter handles deleting a newsletter.
func (h *Handler) DeleteNewsletter(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(idStr); err != nil {
		http.Error(w, `{"error": "invalid newsletter ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteNewsletter(r.Context(), newsletterID); err != nil {
		http.Error(w, `{"error": "failed to delete newsletter"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreatePost handles the creation of a new post.
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	newsletterIDStr := chi.URLParam(r, "id")
	if newsletterIDStr == "" {
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterIDStr); err != nil {
		http.Error(w, `{"error": "invalid newsletter ID"}`, http.StatusBadRequest)
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
	newsletterIDStr := chi.URLParam(r, "id")
	if newsletterIDStr == "" {
		http.Error(w, `{"error": "missing newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterIDStr); err != nil {
		http.Error(w, `{"error": "invalid newsletter ID"}`, http.StatusBadRequest)
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
	newsletterIDStr := chi.URLParam(r, "id")
	postIDStr := chi.URLParam(r, "postID")
	if newsletterIDStr == "" || postIDStr == "" {
		http.Error(w, `{"error": "missing newsletter or post ID"}`, http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterIDStr); err != nil {
		http.Error(w, `{"error": "invalid newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var postID id.Post
	if err := postID.FromString(postIDStr); err != nil {
		http.Error(w, `{"error": "invalid post ID"}`, http.StatusBadRequest)
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
	newsletterIDStr := chi.URLParam(r, "id")
	postIDStr := chi.URLParam(r, "postID")
	if newsletterIDStr == "" || postIDStr == "" {
		http.Error(w, `{"error": "missing newsletter or post ID"}`, http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterIDStr); err != nil {
		http.Error(w, `{"error": "invalid newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var postID id.Post
	if err := postID.FromString(postIDStr); err != nil {
		http.Error(w, `{"error": "invalid post ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePost(r.Context(), newsletterID, postID); err != nil {
		http.Error(w, `{"error": "failed to delete post"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PublishPost(w http.ResponseWriter, r *http.Request) {
	newsletterIDStr := chi.URLParam(r, "id")
	postIDStr := chi.URLParam(r, "postID")
	if newsletterIDStr == "" || postIDStr == "" {
		http.Error(w, `{"error": "missing newsletter or post ID"}`, http.StatusBadRequest)
		return
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterIDStr); err != nil {
		http.Error(w, `{"error": "invalid newsletter ID"}`, http.StatusBadRequest)
		return
	}

	var postID id.Post
	if err := postID.FromString(postIDStr); err != nil {
		http.Error(w, `{"error": "invalid post ID"}`, http.StatusBadRequest)
		return
	}

	p, err := h.service.PublishPost(r.Context(), newsletterID, postID)
	if err != nil {
		http.Error(w, `{"error": "failed to publish post"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
