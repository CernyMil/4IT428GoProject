package api

import (
	"net/http"

	"newsletter-service/service"
	v1 "newsletter-service/transport/api/v1"
	"newsletter-service/transport/middleware"

	firebase "firebase.google.com/go"
	"github.com/go-chi/chi/v5"
)

func NewRouter(newsletterService service.NewsletterService, firebaseApp *firebase.App) http.Handler {
	r := chi.NewRouter()

	authenticator := middleware.NewFirebaseAuthenticator(firebaseApp)
	v1Handler := v1.NewHandler(*authenticator, newsletterService)
	r.Mount("/api/v1", v1Handler)

	return r
}
