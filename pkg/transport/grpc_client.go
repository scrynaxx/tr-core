package transport

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func DefaultDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(forwardHeadersInterceptor()),
	}
}

func forwardHeadersInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		incomingMD, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = metadata.NewOutgoingContext(ctx, incomingMD)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
