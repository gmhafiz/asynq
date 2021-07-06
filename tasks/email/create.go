package email

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"

	"tasks/tasks"
)

type DeliveryPayload struct {
	UserID     int
	TemplateID string
}

type ImageResizePayload struct {
	SourceURL string
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func NewEmailDeliveryTask(userID int, tmplID string) (*asynq.Task, error) {
	// Using msgpack encode the struct to smaller size, thus saving RAM in Redis
	// server.
	// But you cannot inspect a msgpack encoded binary in asynqmon
	payload, err := msgpack.Marshal(&DeliveryPayload{UserID: userID, TemplateID: tmplID})

	// A simple json.Marshal allow json values to be inspected at asynqmon
	// at the cost of bigger RAM usage.
	//payload, err := json.Marshal(DeliveryPayload{UserID: userID, TemplateID: tmplID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(tasks.TypeEmailDelivery, payload), nil
}

func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{SourceURL: src})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(tasks.TypeImageResize, payload), nil
}
