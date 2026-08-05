package employee

import (
	"context"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
	"github.com/scrynaxx/tr-core/internal/modules/employee/repository"
	"github.com/scrynaxx/tr-core/internal/modules/employee/repository/persistence/employee/record"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) repository.Employee {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Get(ctx context.Context, employeeID uuid.UUID) (model.Employee, error) {
	const sql = `
		SELECT *
		FROM employee.employees
		WHERE id = @id`

	rec, err := postgres.Get[record.Employee](ctx, r.pool, sql, pgx.NamedArgs{"id": employeeID}, postgres.WithNotFound(model.ErrEmployeeNotFound))
	if err != nil {
		return model.Employee{}, err
	}

	return record.ToEmployee(rec), nil
}

func (r *Repository) GetByAccount(ctx context.Context, accountID uuid.UUID) (model.Employee, error) {
	const sql = `
		SELECT *
		FROM employee.employees
		WHERE account_id = @account_id`

	rec, err := postgres.Get[record.Employee](ctx, r.pool, sql, pgx.NamedArgs{"account_id": accountID}, postgres.WithNotFound(model.ErrEmployeeNotFound))
	if err != nil {
		return model.Employee{}, err
	}

	return record.ToEmployee(rec), nil
}

func (r *Repository) List(ctx context.Context) ([]model.Employee, error) {
	const sql = `
		SELECT *
		FROM employee.employees
		ORDER BY created_at`

	recs, err := postgres.Select[record.Employee](ctx, r.pool, sql, nil)
	if err != nil {
		return nil, err
	}

	return record.ToListEmployee(recs), nil
}

func (r *Repository) Create(ctx context.Context, employee model.Employee) error {
	const sql = `
		INSERT INTO employee.employees (id, type, first_name, last_name, patronymic, phone, birth_date)
		VALUES (@id, @type, @first_name, @last_name, @patronymic, @phone, @birth_date)`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":         employee.ID,
		"type":       employee.Type,
		"first_name": employee.FirstName,
		"last_name":  employee.LastName,
		"patronymic": employee.Patronymic,
		"phone":      employee.Phone,
		"birth_date": employee.BirthDate,
	})
}

func (r *Repository) Update(ctx context.Context, employee model.Employee) error {
	const sql = `
		UPDATE employee.employees
		SET account_id = @account_id,
		    type = @type,
			first_name = @first_name,
			last_name = @last_name,
			patronymic = @patronymic,
			phone = @phone,
			birth_date = @birth_date,
			archived_at = @archived_at,
			updated_at = now()
		WHERE id = @id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":          employee.ID,
		"account_id":  employee.AccountID,
		"type":        employee.Type,
		"first_name":  employee.FirstName,
		"last_name":   employee.LastName,
		"patronymic":  employee.Patronymic,
		"phone":       employee.Phone,
		"birth_date":  employee.BirthDate,
		"archived_at": employee.ArchivedAt,
	})
}
