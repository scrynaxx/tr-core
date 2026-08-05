package employee

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/internal/modules/auth/client"
	restController "github.com/scrynaxx/tr-core/internal/modules/employee/controller/rest"
	"github.com/scrynaxx/tr-core/internal/modules/employee/migrations"
	employeeRepository "github.com/scrynaxx/tr-core/internal/modules/employee/repository/persistence/employee"
	passportRepository "github.com/scrynaxx/tr-core/internal/modules/employee/repository/persistence/passport"
	employeeUseCase "github.com/scrynaxx/tr-core/internal/modules/employee/usecase/employee"
	passportUseCase "github.com/scrynaxx/tr-core/internal/modules/employee/usecase/passport"
	"github.com/scrynaxx/tr-core/internal/security"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/pkg/messaging"
)

const Schema = "employee"

func New(
	ctx context.Context,
	e *echo.Echo,
	pg *postgres.Postgres,
	bus *messaging.Bus,
	authClient authclient.Client,
	authorizer security.Authorizer,
) error {
	if err := postgres.Migrate(ctx, pg.Pool, Schema, migrations.Files); err != nil {
		return err
	}

	outbox, err := bus.AddOutbox(ctx, pg.Pool, Schema)
	if err != nil {
		return err
	}

	employeeRepo := employeeRepository.New(pg.Pool)
	passportRepo := passportRepository.New(pg.Pool)
	employeeUc := employeeUseCase.New(employeeRepo, passportRepo, outbox, pg.Transactor)
	passportUc := passportUseCase.New(passportRepo)

	restController.RegisterRoutes(e, authorizer, employeeUc, passportUc, authClient)

	return nil
}
