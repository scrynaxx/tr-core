package app

import (
	"context"
	"fmt"
	"time"

	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/pkg/transport/microservice"
	grpcController "github.com/scrynaxx/tr-core/services/order/internal/controller/grpc"
)

const (
	serviceName              = "order"
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

	ms := microservice.New(grpcPort, httpPort, nil)

	if err := grpcController.RegisterRoutes(ms.GRPCServer); err != nil {
		return fmt.Errorf("grpc controller: %w", err)
	}

	return ms.Run(ctx, transportShutdownTimeout)
}
