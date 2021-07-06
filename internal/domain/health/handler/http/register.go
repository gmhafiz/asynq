package http

import (
	"github.com/go-chi/chi/v5"

	"tasks/internal/domain/health"
)

func RegisterHTTPEndPoints(router *chi.Mux, uc health.UseCase) {
	h := NewHandler(uc)

	router.Route("/health", func(router chi.Router) {
		router.Get("/", h.Health)
		router.Get("/readiness", h.Readiness)
	})
}
