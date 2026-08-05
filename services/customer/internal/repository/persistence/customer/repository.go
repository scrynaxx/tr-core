package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/services/customer/internal/model"
	"github.com/scrynaxx/tr-core/services/customer/internal/repository/persistence/customer/record"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Get(ctx context.Context, customerID uuid.UUID) (model.Customer, error) {
	sql := `
	SELECT *
	FROM customer.customers
	WHERE id = @id`

	rec, err := postgres.Get[record.Customer](ctx, r.pool, sql, pgx.NamedArgs{
		"id": customerID,
	}, postgres.WithNotFound(model.ErrCustomerNotFound))
	if err != nil {
		return model.Customer{}, err
	}

	return record.ToCustomer(rec), nil
}

func (r *Repository) List(ctx context.Context) ([]model.Customer, error) {
	sql := `
	SELECT *
	FROM customer.customers
	WHERE archived_at IS NULL
	ORDER BY name, id`

	recs, err := postgres.Select[record.Customer](ctx, r.pool, sql, pgx.NamedArgs{})
	if err != nil {
		return nil, err
	}

	return record.ToCustomers(recs), nil
}

func (r *Repository) Create(ctx context.Context, customer model.Customer) (model.Customer, error) {
	sql := `
	INSERT INTO customer.customers (id, name, phone, email)
	VALUES (@id, @name, @phone, @email)
	RETURNING *`

	rec, err := postgres.Get[record.Customer](ctx, r.pool, sql, pgx.NamedArgs{
		"id":    customer.ID,
		"name":  customer.Name,
		"phone": customer.Phone,
		"email": customer.Email,
	}, postgres.WithExists(model.ErrCustomerExists))
	if err != nil {
		return model.Customer{}, err
	}

	return record.ToCustomer(rec), nil
}

func (r *Repository) Save(ctx context.Context, customer model.Customer) (model.Customer, error) {
	sql := `
	UPDATE customer.customers
	SET name = @name,
		phone = @phone,
		email = @email,
		updated_at = now()
	WHERE id = @id AND archived_at IS NULL
	RETURNING *`

	rec, err := postgres.Get[record.Customer](ctx, r.pool, sql, pgx.NamedArgs{
		"id":    customer.ID,
		"name":  customer.Name,
		"phone": customer.Phone,
		"email": customer.Email,
	}, postgres.WithNotFound(model.ErrCustomerNotFound), postgres.WithExists(model.ErrCustomerExists))
	if err != nil {
		return model.Customer{}, err
	}

	return record.ToCustomer(rec), nil
}

func (r *Repository) Archive(ctx context.Context, customerID uuid.UUID) error {
	sql := `
	UPDATE customer.customers
	SET archived_at = COALESCE(archived_at, now()),
		updated_at = now()
	WHERE id = @id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id": customerID,
	})
}
