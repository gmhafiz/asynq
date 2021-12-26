package email

import (
	"context"
	"encoding/json"
	"fmt"
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

// NewEmailProcessor is constructor that allows us to inject any dependencies
// required by our ProcessTask method.
func NewEmailProcessor(db *sqlx.DB) *Processor {
	// ... return an instance
	return &Processor{
		db: db,
	}
}

// ProcessTask is a method of Processor struct, which implements Handler
// interface.
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

	log.Printf("Task ID: %s Sending Email to User: user_id=%d, From=%d", task.ResultWriter().TaskID(), r.Parameters.To, r.SentByUserID)

	// Email delivery code ...

	// sleep 60 seconds to simulate (hard) work.
	_, err := p.db.ExecContext(ctx, "SELECT SLEEP(60);")
	if err != nil {
		log.Printf("error performing database operation: %v\n", err)
		return fmt.Errorf("database error: %w", err)
	}

	log.Printf("Completed Processing Task ID: %v", task.ResultWriter().TaskID())

	return nil
}
