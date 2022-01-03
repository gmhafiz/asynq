package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"tasks/internal/domain/email"
	"tasks/internal/utility/respond"
	"tasks/internal/utility/validate"
	"tasks/task"
)

type Handler struct {
	validate *validator.Validate
	useCase  email.UseCase
}

// NewHandler receives dependencies and stores them in a local Handler struct.
// All methods of Handler like Send() will be able to use the dependencies.
// Every handler does
//   1. Bind byte stream into local struct
//   2. Validate user input
//   3. Send to appropriate use case layer
//   4. Check response from use case layer
//   5. Respond with result or errors
//
// Note that there are no business logic done in this layer at all.
func NewHandler(uc email.UseCase, validate *validator.Validate) *Handler {
	return &Handler{
		useCase:  uc,
		validate: validate,
	}
}

// Send an email with payload
// @Summary
// @Description
// @Success 200
// @Failure 400
// @Failure 500
// @Router /v1/email [post]
func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {

	// After going through global middlewares, the first thing we want to do is
	// to validate user request.
	// So we make a custom struct with the fields we expect.
	// We need to transform the stream of bytes requests by binding the json
	// input into the struct.
	var req email.Request
	err := email.Bind(r.Body, &req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	errs := validate.Validate(h.validate, req)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	// Depending on the task type constant defined in task/tasks.go, we
	// determine what to do with it.
	switch req.Type {
	case task.TypeEmailDelivery:
		err = h.useCase.Send(r.Context(), req)
		if err != nil {
			respond.Error(w, http.StatusInternalServerError, err)
			return
		}
	default:
		respond.Error(w, http.StatusBadRequest, respond.ErrBadRequest)
	}
}
