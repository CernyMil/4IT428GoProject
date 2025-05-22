package v1

import "github.com/go-chi/chi"

type Handler struct {
	*chi.Mux

	service Service
}

func NewHandler(
	service Service,
) *Handler {
	h := &Handler{
		service: service,
	}
	h.initRouter()
	return h
}

func (h *Handler) initRouter() {
	r := chi.NewRouter()

	// TODO: Setup middleware.

	r.Route("/Editors", func(r chi.Router) {
		r.Get("/", h.ListEditors)
		r.Post("/", h.CreateEditor)
		r.Get("/{email}", h.GetEditor)
		r.Put("/{email}", h.UpdateEditor)
		r.Delete("/{email}", h.DeleteEditor)
	})
	h.Mux = r
}
