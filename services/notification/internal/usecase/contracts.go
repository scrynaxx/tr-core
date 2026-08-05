package usecase

import (
	"context"
)

type Email interface {
	// Send Отправляет письмо получателям из списка receivers
	Send(ctx context.Context, subject, body string, receivers []string) error
}
