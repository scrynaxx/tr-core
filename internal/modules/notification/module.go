package notification

import (
	"context"

	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/internal/modules/notification/repository/webapi/email"
	emailUseCase "github.com/scrynaxx/tr-core/internal/modules/notification/usecase/email"
)

func New(_ context.Context, cfg config.Config) error {
	emailAPI := email.NewAPI(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender, cfg.SMTP.Name)

	_ = emailUseCase.New(emailAPI)

	return nil
}
