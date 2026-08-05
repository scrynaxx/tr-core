package request

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

// HeaderFromContext возвращает первое непустое значение заголовка.
func HeaderFromContext(ctx context.Context, name string) string {
	values := metadata.ValueFromIncomingContext(ctx, strings.ToLower(name))
	for i := range values {
		if values[i] != "" {
			return values[i]
		}
	}

	return ""
}
