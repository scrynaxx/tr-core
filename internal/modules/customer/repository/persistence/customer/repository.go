package customer

import (
	"context"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/internal/modules/customer/model"
	"github.com/scrynaxx/tr-core/internal/modules/customer/repository/persistence/customer/record"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Get(ctx context.Context, customerID uuid.UUID) (model.Customer, error) {
	const sql = `
		SELECT *
		FROM customer.customers
		WHERE id = @id`

	rec, err := postgres.Get[record.Customer](ctx, r.pool, sql, pgx.NamedArgs{"id": customerID}, postgres.WithNotFound(model.ErrCustomerNotFound))
	if err != nil {
		return model.Customer{}, err
	}

	return record.ToCustomer(rec), nil
}

func (r *Repository) List(ctx context.Context) ([]model.Customer, error) {
	const sql = `
		SELECT *
		FROM customer.customers
		ORDER BY created_at`

	recs, err := postgres.Select[record.Customer](ctx, r.pool, sql, pgx.NamedArgs{})
	if err != nil {
		return nil, err
	}

	return record.ToListCustomer(recs), nil
}

func (r *Repository) Create(ctx context.Context, customer model.Customer) error {
	const sql = `
		INSERT INTO customer.customers (id, first_name, last_name, patronymic, phone, email)
		VALUES (@id, @first_name, @last_name, @patronymic, @phone, @email)`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":         customer.ID,
		"first_name": customer.FirstName,
		"last_name":  customer.LastName,
		"patronymic": customer.Patronymic,
		"phone":      customer.Phone,
		"email":      customer.Email,
	}, postgres.WithExists(model.ErrCustomerAlreadyExists))
}

func (r *Repository) Update(ctx context.Context, customer model.Customer) error {
	const sql = `
		UPDATE customer.customers
		SET account_id = @account_id,
		    first_name = @first_name,
		    last_name = @last_name,
		    patronymic = @patronymic,
			phone = @phone,
			email = @email,
			archived_at = @archived_at,
			updated_at = now()
		WHERE id = @id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":          customer.ID,
		"account_id":  customer.AccountID,
		"first_name":  customer.FirstName,
		"last_name":   customer.LastName,
		"patronymic":  customer.Patronymic,
		"phone":       customer.Phone,
		"email":       customer.Email,
		"archived_at": customer.ArchivedAt,
	}, postgres.WithExists(model.ErrCustomerAlreadyExists))
}
