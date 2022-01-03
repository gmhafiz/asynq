package email

import "context"

type UseCase interface {
	Send(context.Context, Request) error
}
