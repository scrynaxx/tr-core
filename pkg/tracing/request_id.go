package tracing

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

const requestIDKey = "request.id"

func WithRequestID(ctx context.Context, id string) context.Context {
	member, err := baggage.NewMember(requestIDKey, id)
	if err != nil {
		return ctx
	}

	bag := baggage.FromContext(ctx)
	bag, err = bag.SetMember(member)
	if err != nil {
		return ctx
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("request.id", id),
	)

	return baggage.ContextWithBaggage(ctx, bag)
}

func GetRequestID(ctx context.Context) string {
	return baggage.FromContext(ctx).Member(requestIDKey).Value()
}
