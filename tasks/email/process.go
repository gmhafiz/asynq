package email

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"

	deliveryV1 "tasks/api/v1"
)

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

// Processor implements asynq.Handler interface.
type Processor struct {
	db *sqlx.DB
}

func NewEmailProcessor(db *sqlx.DB) *Processor {
	// ... return an instance
	return &Processor{
		db: db,
	}
}

func (p Processor) ProcessTask(ctx context.Context, task *asynq.Task) error {
	// example decoding a protobuf encoded payload
	//var p delivery_v1.Delivery
	//err := proto.Unmarshal(t.Payload(), &p)
	//if err != nil {
	//	return err
	//}

	// example decoding a JSON payload
	var payload deliveryV1.Delivery
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return err
	}

	// example decoding a msgpack encoded payload
	//var p DeliveryPayload
	//if err := msgpack.Unmarshal(t.Payload(), &p); err != nil {
	//	return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	//}

	// Email delivery code ...

	log.Printf("Sending Email to User: user_id=%d, From=%d", payload.Email.To, payload.SentBy)

	log.Println(task.ResultWriter().TaskID())

	_, err := p.db.ExecContext(ctx, "SELECT true;")
	if err != nil {
		log.Println("error performing database operation: %w", err)
	}

	return nil
}
