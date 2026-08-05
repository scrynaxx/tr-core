package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/contracts/generated/services/v1/auth"
	"github.com/scrynaxx/tr-core/contracts/generated/services/v1/customer"
	"github.com/scrynaxx/tr-core/contracts/generated/services/v1/employee"
	"github.com/scrynaxx/tr-core/contracts/generated/services/v1/order"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/pkg/transport/gateway"
)

const (
	serviceName     = "gateway"
	gatewayPort     = 8000
	shutdownTimeout = time.Minute
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:                     cfg.Redis.Address,
		Password:                 cfg.Redis.Password,
		MaintNotificationsConfig: &maintnotifications.Config{Mode: "disabled"},
	})
	defer redisClient.Close()

	if err = redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	service := gateway.New(gateway.Params{
		Port:      gatewayPort,
		Redis:     redisClient,
		JWTIssuer: cfg.JWT.Issuer,
		JWTSecret: cfg.JWT.Secret,
		PublicRoutes: []gateway.PublicRoute{
			{Method: http.MethodPost, Path: "/v1/auth/sign-in"},
			{Method: http.MethodPost, Path: "/v1/auth/refresh"},
		},
		Middleware: []gin.HandlerFunc{cors.New(cors.Config{
			MaxAge:           7200,
			AllowCredentials: true,
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Content-Length", "Content-Disposition", "Cache-Control", "X-Request-ID", "X-Device-ID"},
			ExposeHeaders:    []string{"Content-Type", "Content-Length", "Content-Disposition"},
		})},
	})

	if err = service.RegisterEndpoints(ctx, []gateway.Endpoint{
		gateway.NewEndpoint("auth:8000", v1auth.RegisterAuthServiceHandlerFromEndpoint),
		gateway.NewEndpoint("customer:8000", customer.RegisterCustomerServiceHandlerFromEndpoint),
		gateway.NewEndpoint("employee:8000", employee.RegisterEmployeeServiceHandlerFromEndpoint),
		gateway.NewEndpoint("order:8000", v1order.RegisterOrderServiceHandlerFromEndpoint),
	}); err != nil {
		return fmt.Errorf("register endpoints: %w", err)
	}

	return service.Run(ctx, shutdownTimeout)
}
