package grpc

import (
	"github.com/scrynaxx/tr-core/config"
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/auth"
	"github.com/scrynaxx/tr-core/services/auth/internal/controller/grpc/v1"
	"github.com/scrynaxx/tr-core/services/auth/internal/usecase"
	"google.golang.org/grpc"
)

func RegisterRoutes(srv *grpc.Server, environment config.Environment, authCase usecase.Auth) error {
	pbv1.RegisterAuthServiceServer(srv, v1.NewController(environment, authCase))

	return nil
}
