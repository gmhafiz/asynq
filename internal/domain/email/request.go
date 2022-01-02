package email

import (
	"encoding/json"
	"io"

	"tasks/task"
)

// Email is a custom struct concerned with just Email parameters.
type Email struct {
	To uint64 `json:"to"`

	// To is the ID of the user to be sent to
	// Subject is the Email subject
	Subject string `json:"subject" validate:"required"`

	// Content is the text to be sent in the body of an email
	Content string `json:"content" validate:"required"`
}

// RefereeRequest user request is parsed into this struct.
// Other Unit of Work (UoW) will also embed Request struct.
// Other parameters are then defined is their own struct, in this case an Email
// struct
type RefereeRequest struct {

	// The compulsory common Request struct
	task.Request

	// Specific parameters that RefereeRequest needs, in this case, Email struct
	Parameters Email `json:"email" validate:"required"`
}

// Bind is a convenience function to parse incoming user request in JSON format
// to our custom request struct.
func Bind(body io.ReadCloser, b *RefereeRequest) error {
	return json.NewDecoder(body).Decode(b)
}
