package microservice

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"

	"buf.build/go/protovalidate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func grpcErrorInterceptor(mappings map[error]codes.Code) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error("[service] panic recovered",
					slog.Any("panic", recovered),
					slog.String("stack", string(debug.Stack())),
				)

				err = status.Errorf(codes.Internal, "panic occurred: %s", recovered)
				resp = nil
			}
		}()

		if msg, ok := req.(proto.Message); ok {
			if err := protovalidate.Validate(msg); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
			}
		}

		if resp, err = handler(ctx, req); err == nil {
			return resp, nil
		}

		for mappedError, code := range mappings {
			if errors.Is(err, mappedError) {
				if st, detailError := status.New(code, mappedError.Error()).WithDetails(new(errdetails.ErrorInfo)); detailError == nil {
					return nil, st.Err()
				}

				return nil, status.Error(code, mappedError.Error())
			}
		}

		if _, ok := status.FromError(err); ok && status.Code(err) != codes.Unknown {
			return nil, err
		}

		slog.Error("[service] error occurred",
			slog.String("method", info.FullMethod),
			slog.Any("error", err),
		)

		return nil, status.Error(codes.Internal, err.Error())
	}
}
