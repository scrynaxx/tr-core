package employee

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
	"github.com/scrynaxx/tr-core/services/employee/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type traced struct {
	tracer trace.Tracer
	next   usecase.Employee
}

func withTracing(next usecase.Employee) usecase.Employee {
	return &traced{
		tracer: otel.Tracer("employee-usecase-tracer"),
		next:   next,
	}
}

func (t *traced) Create(ctx context.Context, p model.CreateEmployeeParams) (v model.Employee, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Create")
	defer tracing.EndSpan(span, err)

	return t.next.Create(ctx, p)
}

func (t *traced) Get(ctx context.Context, employeeID uuid.UUID) (v model.Employee, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Get")
	defer tracing.EndSpan(span, err)

	return t.next.Get(ctx, employeeID)
}

func (t *traced) List(ctx context.Context) (r []model.Employee, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "List")
	defer tracing.EndSpan(span, err)

	return t.next.List(ctx)
}

func (t *traced) Update(ctx context.Context, employeeID uuid.UUID, p model.UpdateEmployeeParams) (v model.Employee, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Update")
	defer tracing.EndSpan(span, err)

	return t.next.Update(ctx, employeeID, p)
}

func (t *traced) Archive(ctx context.Context, callerID, employeeID uuid.UUID) (err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "Archive")
	defer tracing.EndSpan(span, err)

	return t.next.Archive(ctx, callerID, employeeID)
}

func (t *traced) FindCredentials(ctx context.Context, employeeID uuid.UUID) (v *model.EmployeeCredentials, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "FindCredentials")
	defer tracing.EndSpan(span, err)

	return t.next.FindCredentials(ctx, employeeID)
}

func (t *traced) FindCredentialsByEmail(ctx context.Context, email string) (v uuid.UUID, v2 *model.EmployeeCredentials, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "FindCredentialsByEmail")
	defer tracing.EndSpan(span, err)

	return t.next.FindCredentialsByEmail(ctx, email)
}

func (t *traced) CreateCredentials(ctx context.Context, employeeID uuid.UUID, p model.CreateCredentialsParams) (v model.EmployeeCredentials, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "CreateCredentials")
	defer tracing.EndSpan(span, err)

	return t.next.CreateCredentials(ctx, employeeID, p)
}

func (t *traced) UpdateCredentials(ctx context.Context, employeeID uuid.UUID, email string) (v model.EmployeeCredentials, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "UpdateCredentials")
	defer tracing.EndSpan(span, err)

	return t.next.UpdateCredentials(ctx, employeeID, email)
}

func (t *traced) DeleteCredentials(ctx context.Context, employeeID uuid.UUID) (err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "DeleteCredentials")
	defer tracing.EndSpan(span, err)

	return t.next.DeleteCredentials(ctx, employeeID)
}

func (t *traced) FindPassport(ctx context.Context, employeeID uuid.UUID) (v *model.EmployeePassport, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "FindPassport")
	defer tracing.EndSpan(span, err)

	return t.next.FindPassport(ctx, employeeID)
}

func (t *traced) CreatePassport(ctx context.Context, employeeID uuid.UUID, p model.CreatePassportParams) (v model.EmployeePassport, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "CreatePassport")
	defer tracing.EndSpan(span, err)

	return t.next.CreatePassport(ctx, employeeID, p)
}

func (t *traced) UpdatePassport(ctx context.Context, employeeID uuid.UUID, p model.UpdatePassportParams) (v model.EmployeePassport, err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "UpdatePassport")
	defer tracing.EndSpan(span, err)

	return t.next.UpdatePassport(ctx, employeeID, p)
}

func (t *traced) DeletePassport(ctx context.Context, employeeID uuid.UUID) (err error) {
	ctx, span := tracing.StartSpan(ctx, t.tracer, "DeletePassport")
	defer tracing.EndSpan(span, err)

	return t.next.DeletePassport(ctx, employeeID)
}
