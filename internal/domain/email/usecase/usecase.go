package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"

	"tasks/internal/domain/email"
	emailTask "tasks/tasks/email"
)

type Email struct {
	redis *asynq.Client
	db    *sqlx.DB
}

func New(redis *asynq.Client, db *sqlx.DB) *Email {
	return &Email{
		redis: redis,
		db:    db,
	}
}

func (u *Email) Send(ctx context.Context, req email.RefereeRequest) error {
	task, err := emailTask.NewEmailDeliveryTask(req)
	if err != nil {
		return fmt.Errorf("could not create task: %w", err)
	}

	// we can do any database operations if you wish.
	_, err = u.db.ExecContext(ctx, "SELECT true;")
	if err != nil {
		return fmt.Errorf("error performing database operations: %w", err)
	}

	info, err := u.redis.Enqueue(task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}

	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	return nil
}
