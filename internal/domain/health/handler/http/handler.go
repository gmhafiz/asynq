package http

import (
	"net/http"

	"tasks/internal/domain/health"
	"tasks/internal/utility/respond"
)

type Handler struct {
	useCase health.UseCase
}

func NewHandler(useCase health.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// Health checks if api is up
// @Summary Checks if API is up
// @Description Hits this API to see if API is running in the server
// @Success 200
// @Failure 500
// @Router /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respond.Render(w, http.StatusOK, nil)
}

// Readiness checks if database is alive
// @Summary Checks if both API and Database are up
// @Description Hits this API to see if both API and Database are running in the server
// @Success 200
// @Failure 500
// @Router /health/readiness [get]
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	err := h.useCase.Readiness()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.Render(w, http.StatusOK, nil)
}
