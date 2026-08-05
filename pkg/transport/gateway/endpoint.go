package gateway

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// EndpointRegistrator регистрирует generated gateway handler микросервиса.
type EndpointRegistrator func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

// Endpoint объединяет адрес микросервиса и регистрируемые для него gateway handlers.
type Endpoint struct {
	address     string
	registrator []EndpointRegistrator
}

// NewEndpoint создаёт описание gateway endpoint микросервиса.
func NewEndpoint(address string, registrators ...EndpointRegistrator) Endpoint {
	return Endpoint{
		address:     address,
		registrator: registrators,
	}
}
