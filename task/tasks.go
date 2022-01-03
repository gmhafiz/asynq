package task

// A list of task types constants to be used across all APIs.
// The convention is to separate the namespace by using a colon (:) or dash (-).
// The colon denotes a hierarchy between the namespace and a dash is to
// separate between words.
// e.g.
// 		email:deliver
//		referee:publish
//		export:match-allocations
//
// It is recommended to use a const to define task type to reduce mistakes from
// a typo.
//
// You use the constant in many places, when scheduling a new task, registering
// it in the mux.
const (
	//generate:asynq:gen

	// TypeLongRunningWork simulates a long-running hard work.
	TypeLongRunningWork = "test:heavy-work"

	// TypeEmailDelivery defines email delivery task.
	TypeEmailDelivery = "email:deliver"
)

// Request is a shared struct used by all incoming messages. Typically, you
// want to embed this struct in your own.
//
// For example:
//		type RefereeRequest struct {
//			task.Request
//
//			Parameters Email `json:"email" validate:"required"`
//		}
type Request struct {

	// SentByUserID keeps track of who is sending the message
	SentByUserID uint64 `json:"sent_by" validate:"required"`

	// Type tells what kind of UoW it is. They are defined by constants in
	// task/tasks.go and must be followed by all services using this Producer
	// api.
	Type string `json:"type" validate:"required"`

	// UUID defines a unique ID sent by client. We can send it by header
	// request, or as part oj JSON payload.
	// This field is important to ensure deduplication in the event of client
	// retries.
	UUID string `json:"request_uuid" validate:"uuid4,required"`
}
