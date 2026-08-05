package app

import (
	"context"
	"fmt"
	"time"

	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/pkg/transport/microservice"
	grpcController "github.com/scrynaxx/tr-core/services/customer/internal/controller/grpc"
	"github.com/scrynaxx/tr-core/services/customer/internal/model"
	customerRepository "github.com/scrynaxx/tr-core/services/customer/internal/repository/persistence/customer"
	customerUseCase "github.com/scrynaxx/tr-core/services/customer/internal/usecase/customer"
	"github.com/scrynaxx/tr-core/services/customer/migrations"
	"google.golang.org/grpc/codes"
)

const (
	serviceName              = "customer"
	schemaName               = "customer"
	grpcPort                 = 8000
	httpPort                 = 8001
	transportShutdownTimeout = time.Minute
)

func Run(ctx context.Context, cfg *config.Config) error {
	shutdownTracer, err := tracing.New(ctx, tracing.Config{
		Enabled:     true,
		ServiceName: serviceName,
		Endpoint:    cfg.Tracing.Endpoint,
		Insecure:    true,
		SampleRate:  1.0,
	})
	if err != nil {
		return fmt.Errorf("tracing: %w", err)
	}
	defer shutdownTracer()

	pg, err := postgres.NewPool(ctx, postgres.Params{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Database: cfg.Postgres.Database,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		SSLMode:  cfg.Postgres.SSLMode,
	}, postgres.WithMigrations(schemaName, migrations.Files))
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer pg.Close()

	customerRepo := customerRepository.NewRepository(pg)
	customerCase := customerUseCase.NewUseCase(customerRepo)

	ms := microservice.New(grpcPort, httpPort, map[error]codes.Code{
		model.ErrCustomerNotFound: codes.NotFound,
		model.ErrCustomerExists:   codes.AlreadyExists,
	})

	grpcController.RegisterRoutes(ms.GRPCServer, customerCase)

	return ms.Run(ctx, transportShutdownTimeout)
}
