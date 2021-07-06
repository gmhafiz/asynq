package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"tasks/internal/domain/email"

)

func RegisterHTTPEndPoints(router *chi.Mux, validator *validator.Validate, uc email.UseCase) {
	h := NewHandler(uc, validator)

	router.Route("/api/v1/email", func(router chi.Router) {
		router.Post("/", h.Send)
	})
}
