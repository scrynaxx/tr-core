package grpc

import (
	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/customer"
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/customer"
	"github.com/scrynaxx/tr-core/services/customer/internal/controller/grpc/inner"
	"github.com/scrynaxx/tr-core/services/customer/internal/controller/grpc/v1"
	"github.com/scrynaxx/tr-core/services/customer/internal/usecase"
	"google.golang.org/grpc"
)

func RegisterRoutes(server *grpc.Server, customerCase usecase.Customer) {
	pbinner.RegisterCustomerServiceServer(server, inner.NewController(customerCase))
	pbv1.RegisterCustomerServiceServer(server, v1.NewController(customerCase))
}
