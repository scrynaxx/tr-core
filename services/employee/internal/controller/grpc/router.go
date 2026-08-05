package grpc

import (
	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/employee"
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/employee"
	"github.com/scrynaxx/tr-core/services/employee/internal/controller/grpc/inner"
	"github.com/scrynaxx/tr-core/services/employee/internal/controller/grpc/v1"
	"github.com/scrynaxx/tr-core/services/employee/internal/usecase"
	"google.golang.org/grpc"
)

func RegisterRoutes(server *grpc.Server, employeeCase usecase.Employee) {
	pbinner.RegisterEmployeeServiceServer(server, inner.NewController(employeeCase))
	pbv1.RegisterEmployeeServiceServer(server, v1.NewController(employeeCase))
}
