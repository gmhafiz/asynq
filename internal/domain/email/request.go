package email

import (
	"encoding/json"
	"io"
)

type Request struct {
	Subject      string `json:"subject" validate:"required"`
	Content      string `json:"content" validate:"required"`
	SentByUserID uint64 `json:"sent_by"`
}

type RefereeRequest struct {
	Request
	Parameters Parameter `json:"parameters" validate:"required"`
	Type       string    `json:"type" validate:"required"`
}

type Parameter struct {
	To uint64 `json:"to"`
}

func Bind(body io.ReadCloser, b *RefereeRequest) error {
	return json.NewDecoder(body).Decode(b)
}
