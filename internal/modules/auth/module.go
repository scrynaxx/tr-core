package auth

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/internal/modules/auth/client"
	eventController "github.com/scrynaxx/tr-core/internal/modules/auth/controller/event"
	restController "github.com/scrynaxx/tr-core/internal/modules/auth/controller/rest"
	"github.com/scrynaxx/tr-core/internal/modules/auth/migrations"
	accountRepository "github.com/scrynaxx/tr-core/internal/modules/auth/repository/persistence/account"
	sessionRepository "github.com/scrynaxx/tr-core/internal/modules/auth/repository/persistence/session"
	accountUseCase "github.com/scrynaxx/tr-core/internal/modules/auth/usecase/account"
	authUseCase "github.com/scrynaxx/tr-core/internal/modules/auth/usecase/auth"
	"github.com/scrynaxx/tr-core/internal/security"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"

	"github.com/scrynaxx/tr-core/pkg/messaging"
)

const schema = "auth"

func New(
	ctx context.Context,
	cfg config.Config,
	e *echo.Echo,
	pg *postgres.Postgres,
	bus *messaging.Bus,
	tokenizer security.Tokenizer,
	authorizer security.Authorizer,
) (authclient.Client, error) {
	if err := postgres.Migrate(ctx, pg.Pool, schema, migrations.Files); err != nil {
		return nil, err
	}

	outbox, err := bus.AddOutbox(ctx, pg.Pool, schema)
	if err != nil {
		return nil, err
	}

	accountRepo := accountRepository.New(pg.Pool)
	sessionRepo := sessionRepository.New(pg.Pool)
	accountUc := accountUseCase.New(accountRepo)
	authUc := authUseCase.New(accountRepo, sessionRepo, authorizer, tokenizer, outbox, pg.Transactor)
	client := authclient.New(accountUc)

	restController.RegisterRoutes(e, authorizer, cfg.App.Environment, authUc)

	if err = eventController.RegisterRoutes(bus, authUc); err != nil {
		return nil, err
	}

	return client, nil
}
