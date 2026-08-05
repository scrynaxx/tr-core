package microservice

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Service управляет HTTP и gRPC transport микросервиса и его lifecycle.
type Service struct {
	doneCh           chan struct{}
	errorCh          chan error
	httpPort         int
	grpcPort         int
	cleanupCallbacks []func()
	afterCallbacks   []func()

	HTTPServer *http.Server
	GRPCServer *grpc.Server
	Router     *gin.Engine
}

const (
	httpReadHeaderTimeout = 5 * time.Second
)

// New создаёт transport микросервиса.
func New(grpcPort, httpPort int, errorsMap map[error]codes.Code, options ...grpc.ServerOption) *Service {
	options = append(options,
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(grpcErrorInterceptor(errorsMap)),
	)
	router := newRouter()

	return &Service{
		doneCh:   make(chan struct{}),
		errorCh:  make(chan error, 2),
		httpPort: httpPort,
		grpcPort: grpcPort,
		HTTPServer: &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           router,
			ReadHeaderTimeout: httpReadHeaderTimeout,
		},
		GRPCServer: grpc.NewServer(options...),
		Router:     router,
	}
}

// BeforeTransportClose добавляет обработчики, выполняемые до остановки HTTP и gRPC.
func (m *Service) BeforeTransportClose(callbacks ...func()) {
	m.cleanupCallbacks = append(m.cleanupCallbacks, callbacks...)
}

// AfterTransportClose добавляет обработчики, выполняемые после прекращения новых запросов.
func (m *Service) AfterTransportClose(callbacks ...func()) {
	m.afterCallbacks = append(m.afterCallbacks, callbacks...)
}

// Run запускает HTTP и gRPC transport и блокируется до остановки или ошибки сервера.
func (m *Service) Run(ctx context.Context, shutdownTimeout time.Duration) error {
	runCtx, runCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer runCancel()

	go func() {
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", m.grpcPort))
		if err != nil {
			m.errorCh <- fmt.Errorf("grpc listen on %d: %w", m.grpcPort, err)
			return
		}

		if err := m.GRPCServer.Serve(grpcListener); err != nil {
			m.errorCh <- fmt.Errorf("grpc serve: %w", err)
		}
	}()

	go func() {
		if err := m.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			m.errorCh <- fmt.Errorf("http serve: %w", err)
		}
	}()

	slog.Info("[service] service started",
		slog.Int("grpc_port", m.grpcPort),
		slog.Int("http_port", m.httpPort),
	)

	select {
	case <-runCtx.Done():
		slog.Info("[service] shutdown requested", slog.Any("reason", runCtx.Err()))
		m.shutdown(shutdownTimeout)
		return nil
	case startError := <-m.errorCh:
		slog.Error("[service] server stopping with error", slog.Any("error", startError))
		m.shutdown(shutdownTimeout)
		return startError
	}
}

func (m *Service) shutdown(timeout time.Duration) {
	for _, callback := range m.cleanupCallbacks {
		if callback != nil {
			callback()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func() {
		m.GRPCServer.GracefulStop()
		m.HTTPServer.Shutdown(ctx)
		close(m.doneCh)
	}()

	select {
	case <-m.doneCh:
	case <-ctx.Done():
		m.GRPCServer.Stop()
		m.HTTPServer.Close()
	}

	for _, callback := range m.afterCallbacks {
		if callback != nil {
			callback()
		}
	}

	slog.Info("[service] service stopped")
}
