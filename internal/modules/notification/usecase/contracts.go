package usecase

import (
	"context"
)

type Email interface {
	Send(ctx context.Context, subject, body string, receivers []string) error
}
