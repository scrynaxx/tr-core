package employee

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
	"github.com/scrynaxx/tr-core/services/employee/internal/repository/persistence/employee/record"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Get(ctx context.Context, employeeID uuid.UUID) (model.Employee, error) {
	sql := `
	SELECT *
	FROM employee.employees
	WHERE id = @id`

	rec, err := postgres.Get[record.Employee](ctx, r.pool, sql, pgx.NamedArgs{
		"id": employeeID,
	}, postgres.WithNotFound(model.ErrEmployeeNotFound))
	if err != nil {
		return model.Employee{}, err
	}

	return record.ToEmployee(rec), nil
}

func (r *Repository) List(ctx context.Context) ([]model.Employee, error) {
	sql := `
	SELECT *
	FROM employee.employees
	WHERE archived_at IS NULL
	ORDER BY created_at, id`

	recs, err := postgres.Select[record.Employee](ctx, r.pool, sql, nil)
	if err != nil {
		return nil, err
	}

	return record.ToEmployees(recs), nil
}

func (r *Repository) Create(ctx context.Context, employee model.Employee) (model.Employee, error) {
	sql := `
	INSERT INTO employee.employees (id, type, first_name, last_name, patronymic, phone, birth_date) 
	VALUES (@id, @type, @first_name, @last_name, @patronymic, @phone, @birth_date) 
	RETURNING *`

	rec, err := postgres.Get[record.Employee](ctx, r.pool, sql, pgx.NamedArgs{
		"id":         employee.ID,
		"type":       employee.Type,
		"first_name": employee.FirstName,
		"last_name":  employee.LastName,
		"patronymic": employee.Patronymic,
		"phone":      employee.Phone,
		"birth_date": employee.BirthDate,
	})
	if err != nil {
		return model.Employee{}, err
	}

	return record.ToEmployee(rec), nil
}

func (r *Repository) Save(ctx context.Context, employee model.Employee) (model.Employee, error) {
	sql := `
	UPDATE employee.employees
	SET type = @type,
		first_name = @first_name,
		last_name = @last_name,
		patronymic = @patronymic,
		phone = @phone,
		birth_date = @birth_date,
		updated_at = now()
	WHERE id = @id AND archived_at IS NULL
	RETURNING *`

	rec, err := postgres.Get[record.Employee](ctx, r.pool, sql, pgx.NamedArgs{
		"id":         employee.ID,
		"type":       employee.Type,
		"first_name": employee.FirstName,
		"last_name":  employee.LastName,
		"patronymic": employee.Patronymic,
		"phone":      employee.Phone,
		"birth_date": employee.BirthDate,
	}, postgres.WithNotFound(model.ErrEmployeeNotFound), postgres.WithExists(model.ErrEmployeeExists))
	if err != nil {
		return model.Employee{}, err
	}

	return record.ToEmployee(rec), nil
}

func (r *Repository) Archive(ctx context.Context, employeeID uuid.UUID) error {
	sql := `
	UPDATE employee.employees
	SET archived_at = COALESCE(archived_at, now()),
		updated_at = now()
	WHERE id = @id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id": employeeID,
	})
}
