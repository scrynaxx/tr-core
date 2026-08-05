package order

import (
	"context"

	"github.com/labstack/echo/v5"
	restController "github.com/scrynaxx/tr-core/internal/modules/order/controller/rest"
	"github.com/scrynaxx/tr-core/internal/modules/order/migrations"
	dictionaryRepository "github.com/scrynaxx/tr-core/internal/modules/order/repository/persistence/dictionary"
	dictionaryUseCase "github.com/scrynaxx/tr-core/internal/modules/order/usecase/dictionary"
	"github.com/scrynaxx/tr-core/internal/security"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/pkg/messaging"
)

const Schema = "order"

func New(
	ctx context.Context,
	e *echo.Echo,
	pg *postgres.Postgres,
	bus *messaging.Bus,
	authorizer security.Authorizer,
) error {
	if err := postgres.Migrate(ctx, pg.Pool, Schema, migrations.Files); err != nil {
		return err
	}

	_, err := bus.AddOutbox(ctx, pg.Pool, Schema)
	if err != nil {
		return err
	}

	dictionaryRepo := dictionaryRepository.New(pg.Pool)
	dictionaryUc := dictionaryUseCase.New(dictionaryRepo)

	restController.RegisterRoutes(e, authorizer, dictionaryUc)

	return nil
}
