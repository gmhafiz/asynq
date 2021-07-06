package usecase

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"tasks/internal/domain/email"
	email2 "tasks/tasks/email"
)

type Email struct {
	redis *asynq.Client
}

func New(redis *asynq.Client) *Email {
	return &Email{redis: redis}
}

func (u *Email) Send(ctx context.Context, req email.RefereeRequest) error {
	task, err := email2.NewEmailDeliveryTask(int(req.Request.SentByUserID), "")
	if err != nil {
		return fmt.Errorf("could not create task: %v", err)
	}

	info, err := u.redis.Enqueue(task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %v", err)
	}

	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	return nil
}

