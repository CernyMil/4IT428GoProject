package v1

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	*chi.Mux

	service SubscriberService
}

func NewHandler(service SubscriberService) *Handler {
	h := &Handler{
		service: service,
	}
	h.initRouter()
	return h
}

func (h *Handler) initRouter() {
	r := chi.NewRouter()

	r.Route("/newsletters/{newsletterId}", func(r chi.Router) {
		r.Post("/subscribe", h.SubscribeToNewsletter)
		r.Delete("/unsubscribe", h.UnsubscribeFromNewsletter)
		r.Get("/confirm", h.ConfirmSubscription)
	})
	r.Route("/nginx/newsletters", func(r chi.Router) {
		r.Post("/{newsletterId}/posts/publish", h.SendPublishedPost)
		r.Delete("/{newsletterId}/delete", h.DeleteNewsletter)
		r.Post("/create", h.CreateNewsletter)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Not found: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})

	h.Mux = r
}
