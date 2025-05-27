package v1

import "github.com/go-chi/chi"

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
	r.Route("/nginx/newsletters/{newsletterId}", func(r chi.Router) {
		r.Get("/posts/{postId}/publish", h.SendPublishedPost)
		r.Delete("/delete", h.DeleteNewsletter)
	})

	h.Mux = r
}
