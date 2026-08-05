package passport

import (
	"context"

	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
	"github.com/scrynaxx/tr-core/internal/modules/employee/repository"
	"github.com/scrynaxx/tr-core/internal/modules/employee/repository/persistence/passport/record"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) repository.Passport {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Get(ctx context.Context, employeeID uuid.UUID) (model.Passport, error) {
	const sql = `
		SELECT * FROM employee.passports 
	 	WHERE employee_id = @employee_id`

	rec, err := postgres.Get[record.Passport](ctx, r.pool, sql, pgx.NamedArgs{"employee_id": employeeID}, postgres.WithNotFound(model.ErrPassportNotFound))
	if err != nil {
		return model.Passport{}, err
	}

	return record.ToPassport(rec), nil
}

func (r *Repository) Find(ctx context.Context, employeeID uuid.UUID) (*model.Passport, error) {
	const sql = `
		SELECT * FROM employee.passports 
	 	WHERE employee_id = @employee_id`

	rec, err := postgres.Find[record.Passport](ctx, r.pool, sql, pgx.NamedArgs{"employee_id": employeeID})
	if err != nil {
		return nil, err
	}

	return record.ToPassportPtr(rec), nil
}

func (r *Repository) Create(ctx context.Context, employeeID uuid.UUID, passport model.Passport) error {
	const sql = `
		INSERT INTO employee.passports (employee_id, series, number, issued_by, issued_at, department_code)
		VALUES (@employee_id, @series, @number, @issued_by, @issued_at, @department_code)`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id":     employeeID,
		"series":          passport.Series,
		"number":          passport.Number,
		"issued_by":       passport.IssuedBy,
		"issued_at":       passport.IssuedAt,
		"department_code": passport.DepartmentCode,
	})
}

func (r *Repository) Update(ctx context.Context, employeeID uuid.UUID, passport model.Passport) error {
	const sql = `
		UPDATE employee.passports
		SET series = @series, number = @number, issued_by = @issued_by, issued_at = @issued_at, department_code = @department_code, updated_at = now()
		WHERE employee_id = @employee_id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id":     employeeID,
		"series":          passport.Series,
		"number":          passport.Number,
		"issued_by":       passport.IssuedBy,
		"issued_at":       passport.IssuedAt,
		"department_code": passport.DepartmentCode,
	})
}

func (r *Repository) Delete(ctx context.Context, employeeID uuid.UUID) error {
	const sql = `
		DELETE FROM employee.passports 
	   	WHERE employee_id = @employee_id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{"employee_id": employeeID})
}
