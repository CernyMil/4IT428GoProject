package v1

import (
	"github.com/go-chi/chi"

	"subscriber-service/cmd/api/config"
	"subscriber-service/transport/middleware"
)

type Handler struct {
	*chi.Mux
	Config *config.Config

	service SubscriberService
}

func NewHandler(service SubscriberService, cfg *config.Config) *Handler {
	h := &Handler{
		service: service,
		Config:  cfg,
	}
	h.initRouter()
	return h
}

func (h *Handler) initRouter() {
	r := chi.NewRouter()

	// Public routes
	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/subscribe", h.SubscribeToNewsletter)
		r.Get("/confirm", h.ConfirmSubscription)
		r.Get("/unsubscribe", h.UnsubscribeFromNewsletter)
	})

	// Internal routes with shared middleware
	r.Group(func(r chi.Router) {
		r.Use(middleware.InternalOnlyMiddleware(h.Config.ServiceToken))
		r.Route("/internal", func(r chi.Router) {
			r.Post("/publish-post", h.SendPublishedPost)
			r.Delete("/delete-newsletter", h.DeleteNewsletterSubscriptions)
		})
	})

	h.Mux = r
}
