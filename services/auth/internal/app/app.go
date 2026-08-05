package app

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/contracts/generated/services/inner/employee"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/pkg/events"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/pkg/transport"
	"github.com/scrynaxx/tr-core/pkg/transport/microservice"
	eventController "github.com/scrynaxx/tr-core/services/auth/internal/controller/event"
	grpcController "github.com/scrynaxx/tr-core/services/auth/internal/controller/grpc"
	"github.com/scrynaxx/tr-core/services/auth/internal/model"
	sessionRepository "github.com/scrynaxx/tr-core/services/auth/internal/repository/persistence/session"
	authUseCase "github.com/scrynaxx/tr-core/services/auth/internal/usecase/auth"
	sessionUseCase "github.com/scrynaxx/tr-core/services/auth/internal/usecase/session"
	"github.com/scrynaxx/tr-core/services/auth/migrations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	serviceName              = "auth"
	schemaName               = "auth"
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: "disabled",
		},
	})
	defer rdb.Close()

	if err = rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	bus, err := events.New(ctx, events.Params{
		User:     cfg.Rabbit.User,
		Password: cfg.Rabbit.Password,
		Address:  cfg.Rabbit.Address,
		Vhost:    cfg.Rabbit.Vhost,
	}, pg, schemaName, nil)
	if err != nil {
		return fmt.Errorf("event bus: %w", err)
	}

	employeeConn, _ := grpc.NewClient("employee:8000", transport.DefaultDialOptions()...)
	employeeClient := employee.NewEmployeeServiceClient(employeeConn)
	defer employeeConn.Close()

	sessionRepo := sessionRepository.NewRepository(pg)
	transactor := postgres.NewTransactor(pg)
	revokeRepo := transport.NewRevocationRepository(rdb)

	authCase := authUseCase.NewUseCase(cfg.JWT.Issuer, cfg.JWT.Secret, sessionRepo, employeeClient, bus, transactor)
	sessionCase := sessionUseCase.NewUseCase(revokeRepo)

	ms := microservice.New(grpcPort, httpPort, map[error]codes.Code{
		model.ErrInvalidCredentials: codes.Unauthenticated,
		model.ErrInvalidToken:       codes.PermissionDenied,
		model.ErrSessionExpired:     codes.PermissionDenied,
		model.ErrSessionNotFound:    codes.PermissionDenied,
	})
	ms.BeforeTransportClose(bus.Stop)

	if err = grpcController.RegisterRoutes(ms.GRPCServer, cfg.App.Environment, authCase); err != nil {
		return fmt.Errorf("grpc controller: %w", err)
	}

	if err = eventController.RegisterRoutes(bus, authCase, sessionCase); err != nil {
		return fmt.Errorf("event controller: %w", err)
	}

	bus.Start()

	return ms.Run(ctx, transportShutdownTimeout)
}
