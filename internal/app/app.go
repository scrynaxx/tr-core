package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/internal/echoutil"
	"github.com/scrynaxx/tr-core/internal/echoutil/middleware"
	"github.com/scrynaxx/tr-core/internal/modules/auth"
	"github.com/scrynaxx/tr-core/internal/modules/customer"
	"github.com/scrynaxx/tr-core/internal/modules/employee"
	"github.com/scrynaxx/tr-core/internal/modules/notification"
	"github.com/scrynaxx/tr-core/internal/modules/order"
	"github.com/scrynaxx/tr-core/internal/security"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/pkg/messaging"
	"github.com/scrynaxx/tr-core/pkg/tracing"
)

const (
	address         = "0.0.0.0:8000"
	shutdownTimeout = 45 * time.Second
)

func Run(cfg config.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	tracer, err := tracing.NewTracer(ctx, tracing.Options{
		ServiceName: "server",
		Endpoint:    cfg.Tracing.Endpoint,
		Insecure:    true,
		SampleRate:  1,
	})
	if err != nil {
		return err
	}
	defer tracer.Shutdown(ctx)

	pg, err := postgres.New(ctx, postgres.Options{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Database: cfg.Postgres.Database,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		return err
	}
	defer pg.Pool.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: "disabled",
		},
	})
	defer rdb.Close()

	if err = rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	bus, err := messaging.NewBus(ctx, messaging.Options{
		User:     cfg.Rabbit.User,
		Password: cfg.Rabbit.Password,
		Address:  cfg.Rabbit.Address,
		Vhost:    cfg.Rabbit.Vhost,
	}, nil)
	if err != nil {
		return err
	}

	tokenizer := security.NewTokenizer(cfg.JWT.Issuer, cfg.JWT.Secret)
	authorizer := security.NewAuthorizer(rdb, tokenizer)
	e, esc := buildEchoServer(address, shutdownTimeout)

	authClient, err := auth.New(ctx, cfg, e, pg, bus, tokenizer, authorizer)
	if err != nil {
		return err
	}

	if err = notification.New(ctx, cfg); err != nil {
		return err
	}

	if err = employee.New(ctx, e, pg, bus, authClient, authorizer); err != nil {
		return err
	}

	if _, err = customer.New(ctx, e, pg, authorizer); err != nil {
		return err
	}

	if err = order.New(ctx, e, pg, bus, authorizer); err != nil {
		return err
	}

	if err = esc.Start(ctx, e); err != nil {
		return err
	}

	return shutdown(ctx, bus, tracer)
}

func buildEchoServer(address string, shutdownTimeout time.Duration) (*echo.Echo, echo.StartConfig) {
	e := echo.New()
	e.Validator = echoutil.NewValidator()
	e.Use(
		middleware.Recover(),
		middleware.RequestID(),
		middleware.RequestLogger(),
		middleware.ErrorMapper(map[error]int{
			security.ErrInvalidToken:   http.StatusUnauthorized,
			security.ErrSessionRevoked: http.StatusUnauthorized,
			security.ErrForbidden:      http.StatusForbidden,
		}),
		echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     []string{"http://localhost:3000"},
			ExposeHeaders:    []string{middleware.HeaderXRequestID},
		}),
		echootel.NewMiddleware("server"),
	)
	e.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	return e, echo.StartConfig{
		HideBanner:      true,
		HidePort:        true,
		Address:         address,
		GracefulTimeout: shutdownTimeout,
		ListenerAddrFunc: func(addr net.Addr) {
			slog.Info("[http server] started", slog.String("address", addr.String()))
		},
		OnShutdownError: func(err error) {
			slog.Info("[http server] shutdown error occurred", slog.String("error", err.Error()))
		},
	}
}

func shutdown(ctx context.Context, bus *messaging.Bus, tracer *tracing.Tracer) error {
	select {
	case <-ctx.Done():
		return errors.Join(bus.Shutdown(), tracer.Shutdown(ctx))
	case <-time.After(shutdownTimeout):
		return errors.New("shutdown timed out")
	}
}
