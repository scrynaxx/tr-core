package session

import (
	"context"

	"github.com/scrynaxx/tr-core/pkg/events"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/services/auth/internal/model/event"
	"github.com/scrynaxx/tr-core/services/auth/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type traced struct {
	tracer trace.Tracer
	next   usecase.Session
}

func withTracing(next usecase.Session) usecase.Session {
	return &traced{
		tracer: otel.Tracer("session-usecase-tracer"),
		next:   next,
	}
}

func (t *traced) HandleRevoked(ctx context.Context, e events.Event[event.SessionRevokedDataV1]) (err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "HandleRevoked")
	defer tracing.EndSpan(span, err)

	return t.next.HandleRevoked(ctx, e)
}
