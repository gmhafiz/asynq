package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"

	"tasks/internal/domain/email"
	emailTask "tasks/tasks/email"
)

type Email struct {
	redis *asynq.Client
}

func New(redis *asynq.Client) *Email {
	return &Email{redis: redis}
}

func (u *Email) Send(ctx context.Context, req email.RefereeRequest) error {
	task, err := emailTask.NewEmailDeliveryTask(int(req.Request.SentByUserID), "")
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
