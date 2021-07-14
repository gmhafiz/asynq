package email

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	delivery_v1 "tasks/api/v1"

	"github.com/hibiken/asynq"
)

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	// example decoding a protobuf encoded payload
	var p delivery_v1.Delivery
	err := proto.Unmarshal(t.Payload(), &p)
	if err != nil {
		return err
	}

	// example decoding a JSON payload
	//var p DeliveryPayload
	//if err := json.Unmarshal(t.Payload(), &p); err != nil {

	// example decoding a msgpack encoded payload
	//var p DeliveryPayload
	//if err := msgpack.Unmarshal(t.Payload(), &p); err != nil {
	//	return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	//}

	log.Printf("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	// Email delivery code ...
	return nil
}

// ImageProcessor implements asynq.Handler interface.
type ImageProcessor struct {
	// ... fields for struct
	SourceURL string
}

func (p *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	//var payload ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Resizing image: src=%s", p.SourceURL)
	// Image resizing code ...
	return nil
}

func NewImageProcessor() *ImageProcessor {
	// ... return an instance
	return nil
}
