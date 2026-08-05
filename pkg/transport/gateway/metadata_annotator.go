package gateway

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

func metadataAnnotator(_ context.Context, req *http.Request) metadata.MD {
	md := metadata.New(nil)

	for key, values := range req.Header {
		key = strings.ToLower(key)

		if strings.HasPrefix(key, "x-") {
			for i := range values {
				md.Append(key, values[i])
			}
		}
	}

	return md
}
