package app

import (
	"context"
	"fmt"
	"time"

	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/pkg/events"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/pkg/transport/microservice"
	grpcController "github.com/scrynaxx/tr-core/services/employee/internal/controller/grpc"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
	employeeRepository "github.com/scrynaxx/tr-core/services/employee/internal/repository/persistence/employee"
	credentialsRepository "github.com/scrynaxx/tr-core/services/employee/internal/repository/persistence/employee_credentials"
	passportRepository "github.com/scrynaxx/tr-core/services/employee/internal/repository/persistence/employee_passport"
	employeeUseCase "github.com/scrynaxx/tr-core/services/employee/internal/usecase/employee"
	"github.com/scrynaxx/tr-core/services/employee/migrations"
	"google.golang.org/grpc/codes"
)

const (
	serviceName              = "employee"
	schemaName               = "employee"
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

	bus, err := events.New(ctx, events.Params{
		User:     cfg.Rabbit.User,
		Password: cfg.Rabbit.Password,
		Address:  cfg.Rabbit.Address,
		Vhost:    cfg.Rabbit.Vhost,
	}, pg, schemaName, nil)
	if err != nil {
		return fmt.Errorf("event bus: %w", err)
	}

	transactor := postgres.NewTransactor(pg)
	employeeRepo := employeeRepository.NewRepository(pg)
	credentialsRepo := credentialsRepository.NewRepository(pg)
	passportRepo := passportRepository.NewRepository(pg)
	employeeCase := employeeUseCase.NewUseCase(employeeRepo, credentialsRepo, passportRepo, bus, transactor)

	ms := microservice.New(grpcPort, httpPort, map[error]codes.Code{
		model.ErrEmployeeNotFound:    codes.NotFound,
		model.ErrEmployeeExists:      codes.AlreadyExists,
		model.ErrEmailExists:         codes.AlreadyExists,
		model.ErrPassportNotFound:    codes.NotFound,
		model.ErrCredentialsNotFound: codes.NotFound,
	})

	grpcController.RegisterRoutes(ms.GRPCServer, employeeCase)

	return ms.Run(ctx, transportShutdownTimeout)
}
