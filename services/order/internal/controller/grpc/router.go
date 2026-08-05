package grpc

import (
	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/order"
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/order"
	"github.com/scrynaxx/tr-core/services/order/internal/controller/grpc/inner"
	"github.com/scrynaxx/tr-core/services/order/internal/controller/grpc/v1"
	"google.golang.org/grpc"
)

func RegisterRoutes(srv *grpc.Server) error {
	pbinner.RegisterOrderServiceServer(srv, inner.NewController())
	pbv1.RegisterOrderServiceServer(srv, v1.NewController())
	return nil
}
