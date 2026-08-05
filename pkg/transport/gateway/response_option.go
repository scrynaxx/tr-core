package gateway

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/protobuf/proto"
)

func responseOption(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
	headers := w.Header()

	if cookies := headers["Grpc-Metadata-Set-Cookie"]; len(cookies) > 0 {
		for i := range cookies {
			headers.Add("Set-Cookie", cookies[i])
		}
	}

	for key := range headers {
		if strings.HasPrefix(key, "Grpc-Metadata-") {
			headers.Del(key)
		}
	}

	return nil
}
