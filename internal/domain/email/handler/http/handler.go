package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"tasks/internal/domain/email"
	"tasks/internal/utility/respond"
	"tasks/internal/utility/validate"
)

type Handler struct {
	validate *validator.Validate
	useCase  email.UseCase
}

func NewHandler(uc email.UseCase, validate *validator.Validate) *Handler {
	return &Handler{
		useCase:  uc,
		validate: validate,
	}
}

// Send Send an email with payload
// @Summary
// @Description
// @Success 200
// @Failure 500
// @Router /health [post]
func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	var req email.RefereeRequest
	err := email.Bind(r.Body, &req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	errs := validate.Validate(h.validate, req)
	if errs != nil {
		respond.Error(w, http.StatusBadRequest, errs)
		return
	}

	switch req.Type {
	case "referee":
		err = h.useCase.Send(r.Context(), req)
		if err != nil {
			respond.Error(w, http.StatusInternalServerError, nil)
			return
		}
	default:
		respond.Error(w, http.StatusBadRequest, errs)

	}
}
