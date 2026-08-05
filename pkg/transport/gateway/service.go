package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/scrynaxx/tr-core/pkg/transport"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
)

// Params задаёт транспорт и авторизацию внешнего gateway.
type Params struct {
	Port         int
	Redis        redis.Cmdable
	JWTIssuer    string
	JWTSecret    string
	PublicRoutes []PublicRoute
	Middleware   []gin.HandlerFunc
}

// Service обслуживает внешний HTTP API через grpc-gateway.
type Service struct {
	port             int
	server           *http.Server
	mux              *runtime.ServeMux
	errorCh          chan error
	cleanupCallbacks []func()
}

// New создаёт gateway с общей цепочкой авторизации.
func New(params Params) *Service {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaler),
		runtime.WithForwardResponseOption(responseOption),
		runtime.WithMetadata(metadataAnnotator),
		runtime.WithErrorHandler(errorHandler),
	)
	revokeRepo := transport.NewRevocationRepository(params.Redis)
	verifier := newTokenVerifier(params.JWTIssuer, params.JWTSecret)
	router := newRouter(mux, revokeRepo, verifier, params.PublicRoutes, params.Middleware...)

	return &Service{
		port:    params.Port,
		mux:     mux,
		errorCh: make(chan error, 1),
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", params.Port),
			Handler:           router,
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
		},
	}
}

// RegisterEndpoints регистрирует handlers всех downstream микросервисов.
func (g *Service) RegisterEndpoints(ctx context.Context, endpoints []Endpoint) error {
	for _, endpoint := range endpoints {
		for _, registrator := range endpoint.registrator {
			if err := registrator(ctx, g.mux, endpoint.address, transport.DefaultDialOptions()); err != nil {
				return fmt.Errorf("register %s handlers: %w", endpoint.address, err)
			}
		}
	}

	return nil
}

// BeforeTransportShutdown добавляет обработчики, выполняемые перед остановкой HTTP transport.
func (g *Service) BeforeTransportShutdown(callbacks ...func()) {
	g.cleanupCallbacks = append(g.cleanupCallbacks, callbacks...)
}

// Run запускает gateway и блокируется до остановки или ошибки сервера.
func (g *Service) Run(ctx context.Context, shutdownTimeout time.Duration) error {
	runCtx, runCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer runCancel()

	go func() {
		if err := g.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			g.errorCh <- fmt.Errorf("http serve on %d: %w", g.port, err)
		}
	}()

	slog.Info("[gateway] service started", slog.Int("http_port", g.port))

	select {
	case <-runCtx.Done():
		slog.Info("[gateway] service shutdown requested", slog.Any("reason", runCtx.Err()))
		g.shutdown(shutdownTimeout)
		return nil
	case startError := <-g.errorCh:
		slog.Error("[gateway] server stopping with error", slog.Any("error", startError))
		g.shutdown(shutdownTimeout)
		return startError
	}
}

func (g *Service) shutdown(timeout time.Duration) {
	for _, callback := range g.cleanupCallbacks {
		if callback != nil {
			callback()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := g.server.Shutdown(ctx); err != nil {
		_ = g.server.Close()
	}

	slog.Info("[gateway] service stopped")
}
