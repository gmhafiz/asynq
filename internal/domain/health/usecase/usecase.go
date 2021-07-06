package usecase

import "tasks/internal/domain/health"

type Health struct {
	healthRepo health.Repository
}

func New(health health.Repository) *Health {
	return &Health{
		healthRepo: health,
	}
}

func (u *Health) Readiness() error {
	return u.healthRepo.Readiness()
}
