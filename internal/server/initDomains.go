package server

import (
	emailHandler "tasks/internal/domain/email/handler/http"
	emailUseCase "tasks/internal/domain/email/usecase"
	healthHandler "tasks/internal/domain/health/handler/http"
	healthRepo "tasks/internal/domain/health/repository/postgres"
	healthUseCase "tasks/internal/domain/health/usecase"
)

// initDomains groups your tasks into domains.
func (s *Server) initDomains() {
	s.initHealth()
	s.initEmail()
}

// initHealth is a method of Server struct.
// It tends to be simple and the only dependency is a database connection pool
// where it checks whether  this API still has a connection to it.
func (s *Server) initHealth() {
	newHealthRepo := healthRepo.New(s.DB())
	newHealthUseCase := healthUseCase.New(newHealthRepo)
	healthHandler.RegisterHTTPEndPoints(s.router, newHealthUseCase)
}

// initEmail is a method of Server struct.
// This is where dependency injection can happen. Any dependencies you need is
// simply passed into a function as a parameter.
func (s *Server) initEmail() {
	newEmailUseCase := emailUseCase.New(s.asynq, s.DB())
	emailHandler.RegisterHTTPEndPoints(s.router, s.validator, newEmailUseCase)
}
