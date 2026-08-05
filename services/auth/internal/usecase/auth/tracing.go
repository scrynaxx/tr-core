package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/contracts/globalevent"
	"github.com/scrynaxx/tr-core/pkg/events"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/services/auth/internal/model"
	"github.com/scrynaxx/tr-core/services/auth/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type traced struct {
	tracer trace.Tracer
	next   usecase.Auth
}

func withTracing(next usecase.Auth) usecase.Auth {
	return &traced{
		tracer: otel.Tracer("auth-usecase-tracer"),
		next:   next,
	}
}

func (t *traced) SignIn(ctx context.Context, email, password, refresh string) (v *model.TokenData, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "SignIn")
	defer tracing.EndSpan(span, err)

	return t.next.SignIn(ctx, email, password, refresh)
}

func (t *traced) SignOut(ctx context.Context, employeeID, sessionID uuid.UUID) (err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "SignOut")
	defer tracing.EndSpan(span, err)

	return t.next.SignOut(ctx, employeeID, sessionID)
}

func (t *traced) Refresh(ctx context.Context, refresh string, userAgent string) (v *model.TokenData, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Refresh")
	defer tracing.EndSpan(span, err)

	return t.next.Refresh(ctx, refresh, userAgent)
}

func (t *traced) HandleEmployeeArchived(ctx context.Context, e events.Event[globalevent.EmployeeArchivedDataV1]) (err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "HandleEmployeeArchived")
	defer tracing.EndSpan(span, err)

	return t.next.HandleEmployeeArchived(ctx, e)
}

func (t *traced) HandleEmployeeCredentialsChanged(ctx context.Context, e events.Event[globalevent.EmployeeCredentialsChangedDataV1]) (err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "HandleEmployeeCredentialsChanged")
	defer tracing.EndSpan(span, err)

	return t.next.HandleEmployeeCredentialsChanged(ctx, e)
}
