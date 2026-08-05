package email

import (
	"context"
	"fmt"

	"github.com/scrynaxx/tr-core/services/notification/internal/repository"
	"github.com/scrynaxx/tr-core/services/notification/internal/usecase"
)

type UseCase struct {
	emailAPI repository.EmailWebAPI
}

func NewUseCase(emailAPI repository.EmailWebAPI) usecase.Email {
	return withTracing(&UseCase{
		emailAPI: emailAPI,
	})
}

func (s *UseCase) Send(ctx context.Context, subject, body string, receivers []string) error {
	if err := s.emailAPI.Send(ctx, subject, body, receivers); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
