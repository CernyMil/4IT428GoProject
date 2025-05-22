package api

import (
	"net/http"

	"your_project_path/service/newsletter"
	v1 "your_project_path/transport/api/v1"

	"github.com/go-chi/chi/v5"
)

func NewRouter(newsletterService newsletter.Service) http.Handler {
	r := chi.NewRouter()

	v1Handler := v1.NewNewsletterHandler(newsletterService)
	r.Mount("/api/v1", v1Handler.Routes())

	return r
}
