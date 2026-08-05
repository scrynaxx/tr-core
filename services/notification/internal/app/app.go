package app

import (
	"context"
	"fmt"
	"time"

	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/pkg/transport/microservice"
	grpcController "github.com/scrynaxx/tr-core/services/notification/internal/controller/grpc"
	emailWebAPI "github.com/scrynaxx/tr-core/services/notification/internal/repository/webapi/email"
	emailUseCase "github.com/scrynaxx/tr-core/services/notification/internal/usecase/email"
)

const (
	serviceName              = "notification"
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

	emailAPI, err := emailWebAPI.NewAPI(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender, cfg.SMTP.Name)
	if err != nil {
		return fmt.Errorf("email api: %w", err)
	}

	email := emailUseCase.NewUseCase(emailAPI)

	ms := microservice.New(grpcPort, httpPort, nil)

	if err = grpcController.RegisterRoutes(ms.GRPCServer, email); err != nil {
		return fmt.Errorf("grpc controller: %w", err)
	}

	return ms.Run(ctx, transportShutdownTimeout)
}
