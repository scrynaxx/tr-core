package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const defaultShutdownTimeout = 5 * time.Second

type Config struct {
	Enabled     bool
	ServiceName string
	Endpoint    string
	Insecure    bool
	SampleRate  float64
}

func New(ctx context.Context, cfg Config) (func() error, error) {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	if !cfg.Enabled {
		return func() error { return nil }, nil
	}

	exporterOpts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(cfg.Endpoint)}
	if cfg.Insecure {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, exporterOpts...)
	if err != nil {
		return nil, fmt.Errorf("new grpc trace: %w", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("new resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SampleRate))),
	)
	otel.SetTracerProvider(tp)

	return func() error {
		sdCtx, cancel := context.WithTimeout(ctx, defaultShutdownTimeout)
		defer cancel()

		if err := tp.Shutdown(sdCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}

		return nil
	}, nil
}

func EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()
}

func StartSpan(ctx context.Context, tracer trace.Tracer, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, name, opts...)

	if requestID := GetRequestID(ctx); requestID != "" {
		span.SetAttributes(attribute.String("request.id", requestID))
	}

	return ctx, span
}
