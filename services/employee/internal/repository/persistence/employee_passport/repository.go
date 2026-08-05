package employee

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
	"github.com/scrynaxx/tr-core/services/employee/internal/repository/persistence/employee_passport/record"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Get(ctx context.Context, employeeID uuid.UUID) (model.EmployeePassport, error) {
	sql := `
	SELECT *
	FROM employee.passports
	WHERE employee_id = @employee_id`

	rec, err := postgres.Get[record.Passport](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id": employeeID,
	}, postgres.WithNotFound(model.ErrPassportNotFound))

	return record.ToPassport(rec), err
}

func (r *Repository) Find(ctx context.Context, employeeID uuid.UUID) (*model.EmployeePassport, error) {
	sql := `
	SELECT *
	FROM employee.passports
	WHERE employee_id = @employee_id`

	rec, err := postgres.Find[record.Passport](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id": employeeID,
	})

	return record.ToPassportPtr(rec), err
}

func (r *Repository) Create(ctx context.Context, employeeID uuid.UUID, passport model.EmployeePassport) (model.EmployeePassport, error) {
	sql := `
	INSERT INTO employee.passports (employee_id, series, number, issued_by, issued_at, department_code) 
	VALUES (@employee_id, @series, @number, @issued_by, @issued_at, @department_code) 
	RETURNING *`

	rec, err := postgres.Get[record.Passport](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id":     employeeID,
		"series":          passport.Series,
		"number":          passport.Number,
		"issued_by":       passport.IssuedBy,
		"issued_at":       passport.IssuedAt,
		"department_code": passport.DepartmentCode,
	}, postgres.WithNotFound(model.ErrPassportNotFound))
	if err != nil {
		return model.EmployeePassport{}, err
	}

	return record.ToPassport(rec), nil
}

func (r *Repository) Save(ctx context.Context, employeeID uuid.UUID, passport model.EmployeePassport) (model.EmployeePassport, error) {
	sql := `
	UPDATE employee.passports
	SET series = @series,
		number = @number,
		issued_by = @issued_by,
		issued_at = @issued_at,
		department_code = @department_code,
		updated_at = now()
	WHERE employee_id = @employee_id
	RETURNING *`

	rec, err := postgres.Get[record.Passport](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id":     employeeID,
		"series":          passport.Series,
		"number":          passport.Number,
		"issued_by":       passport.IssuedBy,
		"issued_at":       passport.IssuedAt,
		"department_code": passport.DepartmentCode,
	}, postgres.WithNotFound(model.ErrPassportNotFound))
	if err != nil {
		return model.EmployeePassport{}, err
	}

	return record.ToPassport(rec), nil
}

func (r *Repository) Delete(ctx context.Context, employeeID uuid.UUID) error {
	sql := `
	DELETE FROM employee.passports
	WHERE employee_id = @employee_id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id": employeeID,
	})
}
