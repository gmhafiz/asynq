package email

import (
	"encoding/json"
	//"encoding/json"
	"tasks/tasks"

	"github.com/hibiken/asynq"
	//deliveryV1 "tasks/api/v1"
	"tasks/internal/domain/email"
)

//type DeliveryPayload struct {
//	UserID     int
//	TemplateID string
//}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

// NewEmailDeliveryTask user request DTO is converted to protobuf-generated
// struct and then queued to Redis.
// Here are shown 3 ways you can serialize the payload:
//    1. msgpack
//    2. json
//    3. protobuf
func NewEmailDeliveryTask(req email.RefereeRequest) (*asynq.Task, error) {
	// Using msgpack encode the struct to smaller size, thus saving RAM in Redis
	// server.
	// But you cannot inspect a msgpack encoded binary in asynqmon
	//payload, err := msgpack.Marshal(req)
	//if err != nil {
	//	return nil, err
	//}

	// Using protocol buffer encode the struct to an even smaller size than
	// msgpack, thus saving RAM in Redis server on top of being very fast.
	// But you cannot inspect a protobuf encoded binary in asynqmon
	//msg := &deliveryV1.Delivery{
	//	SentBy: int64(req.SentByUserID),
	//	Type:   tasks.TypeEmailDelivery,
	//	Email: &deliveryV1.Email{
	//		Subject: req.Parameters.Subject,
	//		Content: req.Parameters.Content,
	//		To:      int64(req.Parameters.To),
	//	},
	//}
	// Best serialization speed and the smallest payload.
	//payload, err := proto.Marshal(msg)
	//if err != nil {
	//	return nil, err
	//}

	// A simple json.Marshal allow json values to be inspected at asynqmon
	// at the cost of bigger RAM usage and slower serialization.
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(tasks.TypeEmailDelivery, payload), nil
}
