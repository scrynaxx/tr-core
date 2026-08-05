package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

var marshaler = &runtime.HTTPBodyMarshaler{
	Marshaler: &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	},
}
