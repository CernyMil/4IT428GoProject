package api

import (
	_ "embed"
	"log/slog"
	"net/http"

	"subscriber-service/cmd/api/config"
	apiv1 "subscriber-service/transport/api/v1"

	"github.com/go-chi/chi"
	httpx "go.strv.io/net/http"
)

//go:embed openapi.yaml
var OpenAPI []byte

// Controller handles all /api endpoints.
// It is responsible for routing requests to appropriate handlers.
// Versioned endpoints are handled by subcontrollers.
type Controller struct {
	*chi.Mux

	service apiv1.SubscriberService
	cfg     *config.Config
	version string
}

func NewController(service apiv1.SubscriberService, cfg *config.Config, version string) (*Controller, error) {
	controller := &Controller{
		service: service,
		version: version,
		cfg:     cfg,
	}
	controller.initRouter()
	return controller, nil
}

func (c *Controller) initRouter() {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		// r.Use(httpx.LoggingMiddleware(util.NewServerLogger("httpx.LoggingMiddleware")))
		// r.Use(httpx.RecoverMiddleware(util.NewServerLogger("httpx.RecoverMiddleware").WithStackTrace(slog.Level)))

		v1Handler := apiv1.NewHandler(
			c.service,
			c.cfg,
		)

		r.Route("/api", func(r chi.Router) {
			r.Get("/openapi.yaml", c.OpenAPI)
			r.Mount("/v1", v1Handler)
		})
	})

	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Get("/version", c.Version)

	c.Mux = r
}

func (c *Controller) OpenAPI(w http.ResponseWriter, _ *http.Request) {
	if err := httpx.WriteResponse(w, OpenAPI, http.StatusOK); err != nil {
		slog.Error("writing response", slog.Any("error", err))
	}
}

func (c *Controller) Version(w http.ResponseWriter, _ *http.Request) {
	if err := httpx.WriteResponse(w, c.version, http.StatusOK); err != nil {
		slog.Error("writing response", slog.Any("error", err))
	}
}
