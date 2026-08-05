package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/services/customer/internal/model"
	"github.com/scrynaxx/tr-core/services/customer/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type traced struct {
	tracer trace.Tracer
	next   usecase.Customer
}

func withTracing(next usecase.Customer) usecase.Customer {
	return &traced{
		tracer: otel.Tracer("customer-usecase-tracer"),
		next:   next,
	}
}

func (t *traced) Create(ctx context.Context, params model.CreateCustomerParams) (v model.Customer, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Create")
	defer tracing.EndSpan(span, err)

	return t.next.Create(ctx, params)
}

func (t *traced) Get(ctx context.Context, customerID uuid.UUID) (v model.Customer, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Get")
	defer tracing.EndSpan(span, err)

	return t.next.Get(ctx, customerID)
}

func (t *traced) List(ctx context.Context) (v []model.Customer, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "List")
	defer tracing.EndSpan(span, err)

	return t.next.List(ctx)
}

func (t *traced) Update(ctx context.Context, customerID uuid.UUID, params model.UpdateCustomerParams) (v model.Customer, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Save")
	defer tracing.EndSpan(span, err)

	return t.next.Update(ctx, customerID, params)
}

func (t *traced) Archive(ctx context.Context, customerID uuid.UUID) (err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Archive")
	defer tracing.EndSpan(span, err)

	return t.next.Archive(ctx, customerID)
}
