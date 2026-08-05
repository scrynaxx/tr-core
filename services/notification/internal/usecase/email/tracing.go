package email

import (
	"context"

	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/services/notification/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type traced struct {
	tracer trace.Tracer
	next   usecase.Email
}

func withTracing(next usecase.Email) usecase.Email {
	return &traced{
		tracer: otel.Tracer("email-usecase-tracer"),
		next:   next,
	}
}

func (t *traced) Send(ctx context.Context, subject, body string, receivers []string) (err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Send")
	defer tracing.EndSpan(span, err)

	return t.next.Send(ctx, subject, body, receivers)
}
