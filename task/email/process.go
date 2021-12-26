package email

import (
	"context"
	"encoding/json"
	"tasks/internal/domain/email"

	"log"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	//deliveryV1 "tasks/api/v1"
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
	var r email.RefereeRequest
	if err := json.Unmarshal(task.Payload(), &r); err != nil {
		return err
	}

	// example decoding a msgpack encoded payload
	//var r email.RefereeRequest
	//if err := msgpack.Unmarshal(task.Payload(), &r); err != nil {
	//	return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	//}

	// Email delivery code ...

	log.Printf("Sending Email to User: user_id=%d, From=%d", r.Parameters.To, r.SentByUserID)

	log.Println(task.ResultWriter().TaskID())

	_, err := p.db.ExecContext(ctx, "SELECT true;")
	if err != nil {
		log.Println("error performing database operation: %w", err)
	}

	log.Printf("Completed Processing Task ID: %v", task.ResultWriter().TaskID())

	return nil
}
