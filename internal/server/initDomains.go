package server

import (
	emailHandler "tasks/internal/domain/email/handler/http"
	emailUseCase "tasks/internal/domain/email/usecase"
	healthHandler "tasks/internal/domain/health/handler/http"
	healthRepo "tasks/internal/domain/health/repository/postgres"
	healthUseCase "tasks/internal/domain/health/usecase"
)

func (s *Server) initHealth() {
	newHealthRepo := healthRepo.New(s.DB())
	newHealthUseCase := healthUseCase.New(newHealthRepo)
	healthHandler.RegisterHTTPEndPoints(s.router, newHealthUseCase)
}

func (s *Server) initEmail() {
	newEmailUseCase := emailUseCase.New(s.redis)
	emailHandler.RegisterHTTPEndPoints(s.router, s.validator, newEmailUseCase)
}
