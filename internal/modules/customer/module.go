package customer

import (
	"context"

	"github.com/labstack/echo/v5"
	customerclient "github.com/scrynaxx/tr-core/internal/modules/customer/client"
	restController "github.com/scrynaxx/tr-core/internal/modules/customer/controller/rest"
	"github.com/scrynaxx/tr-core/internal/modules/customer/migrations"
	customerRepository "github.com/scrynaxx/tr-core/internal/modules/customer/repository/persistence/customer"
	customerUseCase "github.com/scrynaxx/tr-core/internal/modules/customer/usecase/customer"
	"github.com/scrynaxx/tr-core/internal/security"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
)

const schema = "customer"

func New(
	ctx context.Context,
	e *echo.Echo,
	pg *postgres.Postgres,
	authorizer security.Authorizer,
) (customerclient.Client, error) {
	if err := postgres.Migrate(ctx, pg.Pool, schema, migrations.Files); err != nil {
		return nil, err
	}

	customerRepo := customerRepository.NewRepository(pg.Pool)
	customerCase := customerUseCase.NewUseCase(customerRepo)
	client := customerclient.New()

	restController.RegisterRoutes(e, authorizer, customerCase)

	return client, nil
}
