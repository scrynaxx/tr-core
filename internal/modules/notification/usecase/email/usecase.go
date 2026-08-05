package email

import (
	"context"

	"github.com/scrynaxx/tr-core/internal/modules/notification/repository"
	"github.com/scrynaxx/tr-core/internal/modules/notification/usecase"
)

type UseCase struct {
	emailAPI repository.EmailWebAPI
}

func New(emailAPI repository.EmailWebAPI) usecase.Email {
	return &UseCase{
		emailAPI: emailAPI,
	}
}

func (uc *UseCase) Send(ctx context.Context, subject, body string, receivers []string) error {
	return uc.emailAPI.Send(ctx, subject, body, receivers)
}
