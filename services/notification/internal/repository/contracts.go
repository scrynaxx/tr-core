package repository

import (
	"context"
)

type EmailWebAPI interface {
	Send(ctx context.Context, subject, body string, receivers []string) error
}
