package employee

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
	"github.com/scrynaxx/tr-core/services/employee/internal/repository/persistence/employee_credentials/record"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Get(ctx context.Context, employeeID uuid.UUID) (model.EmployeeCredentials, error) {
	sql := `
	SELECT *
	FROM employee.credentials
	WHERE employee_id = @employee_id`

	rec, err := postgres.Get[record.Credentials](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id": employeeID,
	}, postgres.WithNotFound(model.ErrCredentialsNotFound))
	if err != nil {
		return model.EmployeeCredentials{}, err
	}

	return record.ToCredentials(rec), err
}

func (r *Repository) Find(ctx context.Context, employeeID uuid.UUID) (*model.EmployeeCredentials, error) {
	sql := `
	SELECT *
	FROM employee.credentials
	WHERE employee_id = @employee_id`

	rec, err := postgres.Find[record.Credentials](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id": employeeID,
	})

	return record.ToCredentialsPtr(rec), err
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (uuid.UUID, *model.EmployeeCredentials, error) {
	sql := `
	SELECT *
	FROM employee.credentials
	WHERE email = @email`

	rec, err := postgres.Find[record.Credentials](ctx, r.pool, sql, pgx.NamedArgs{
		"email": email,
	})
	if rec == nil {
		return uuid.Nil, nil, nil
	}

	return rec.EmployeeID, record.ToCredentialsPtr(rec), err
}

func (r *Repository) Create(ctx context.Context, employeeID uuid.UUID, credentials model.EmployeeCredentials) (model.EmployeeCredentials, error) {
	sql := `
	INSERT INTO employee.credentials (employee_id, email, password_hash) 
	VALUES (@employee_id, @email, @password_hash)
	RETURNING *`

	rec, err := postgres.Get[record.Credentials](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id":   employeeID,
		"email":         credentials.Email,
		"password_hash": credentials.PasswordHash,
	}, postgres.WithExists(model.ErrEmailExists))
	if err != nil {
		return model.EmployeeCredentials{}, err
	}

	return record.ToCredentials(rec), nil
}

func (r *Repository) Save(ctx context.Context, employeeID uuid.UUID, credentials model.EmployeeCredentials) (model.EmployeeCredentials, error) {
	sql := `
	UPDATE employee.credentials
	SET email = @email,
		password_hash = @password_hash,
		updated_at = now()
	WHERE employee_id = @employee_id
	RETURNING *`

	rec, err := postgres.Get[record.Credentials](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id":   employeeID,
		"email":         credentials.Email,
		"password_hash": credentials.PasswordHash,
	}, postgres.WithNotFound(model.ErrCredentialsNotFound), postgres.WithExists(model.ErrEmailExists))
	if err != nil {
		return model.EmployeeCredentials{}, err
	}

	return record.ToCredentials(rec), nil
}

func (r *Repository) Delete(ctx context.Context, employeeID uuid.UUID) error {
	sql := `
	DELETE FROM employee.credentials
	WHERE employee_id = @employee_id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id": employeeID,
	})
}
