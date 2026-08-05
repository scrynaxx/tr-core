package grpc

import (
	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/notification"
	"github.com/scrynaxx/tr-core/services/notification/internal/controller/grpc/inner"
	"github.com/scrynaxx/tr-core/services/notification/internal/usecase"
	"google.golang.org/grpc"
)

func RegisterRoutes(srv *grpc.Server, email usecase.Email) error {
	pbinner.RegisterNotificationServiceServer(srv, inner.NewController(email))

	return nil
}
