package usecase

import (
	"context"
	"fmt"
	"log"
	"tasks/internal/utility/queue"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"

	"tasks/internal/domain/email"
	emailTask "tasks/task/email"
)

type Email struct {
	asynqClient *asynq.Client
	db          *sqlx.DB
}

// New is the Use Case layer where all business logic are implemented.
func New(client *asynq.Client, db *sqlx.DB) *Email {
	return &Email{
		asynqClient: client,
		db:          db,
	}
}

func (u *Email) Send(ctx context.Context, req email.Request) error {
	log.Printf("processing: %s\n", req.UUID)

	task, err := emailTask.NewEmailDeliveryTask(ctx, req)
	if err != nil {
		return fmt.Errorf("could not create task: %w", err)
	}

	// We can do any database operations if you wish, preferably making a
	// data access layer (DAL) for it.
	//
	// It is possible to wrap both database query and task enqueue in a
	// database transaction, depending on your use case. This ensures no
	// record is written to database when enqueue has failed.
	_, err = u.db.ExecContext(ctx, "SELECT sleep(20);")
	if err != nil {
		return fmt.Errorf("error performing database operations: %w", err)
	}

	// This is where it enqueues the task to Redis where a consumer will pick
	// up.
	// Enqueue()  accepts a variadic third parameter. You can give it a
	// timeout, retry, schedule it to process after a certain duration, unique
	// or schedule it to send at a certain time.
	// Unique:
	//		https://github.com/hibiken/asynq/wiki/Unique-Tasks
	// The TaskID() option ensures idempotency by allowing only UoW identified
	// by its UUID. It is tracked by redis, so we already achieved a distributed
	// lock on the UUID key.
	// Retention() keeps the UUID in Redis thus preventing the same request with
	// the same UUID to be run again.
	info, err := queue.Enqueue(ctx, u.asynqClient, task, req.UUID)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}

	log.Printf("enqueued task: id=%s queue=%s type=%s", info.ID, info.Queue, info.Type)

	return nil
}
