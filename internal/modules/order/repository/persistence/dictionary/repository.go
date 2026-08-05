package dictionary

import (
	"context"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/internal/modules/order/model"
	"github.com/scrynaxx/tr-core/internal/modules/order/repository"
	"github.com/scrynaxx/tr-core/internal/modules/order/repository/persistence/dictionary/record"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) repository.Dictionary {
	return &Repository{pool: pool}
}

func (r *Repository) GetCargoPackage(ctx context.Context, packageID uuid.UUID) (model.CargoPackage, error) {
	const sql = `
		SELECT *
		FROM order.cargo_packages
		WHERE id = @id`

	rec, err := postgres.Get[record.CargoPackage](ctx, r.pool, sql, pgx.NamedArgs{"id": packageID}, postgres.WithNotFound(model.ErrCargoPackageNotFound))
	if err != nil {
		return model.CargoPackage{}, err
	}

	return record.ToCargoPackage(rec), nil
}

func (r *Repository) ListCargoPackages(ctx context.Context) ([]model.CargoPackage, error) {
	const sql = `
		SELECT *
		FROM order.cargo_packages`

	rec, err := postgres.Select[record.CargoPackage](ctx, r.pool, sql, nil)
	if err != nil {
		return nil, err
	}

	return record.ToListCargoPackages(rec), nil
}

func (r *Repository) CreateCargoPackage(ctx context.Context, pack model.CargoPackage) error {
	const sql = `
		INSERT INTO order.cargo_packages (id, name)
		VALUES (@id, @name)`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":   pack.ID,
		"name": pack.Name,
	})
}

func (r *Repository) UpdateCargoPackage(ctx context.Context, pack model.CargoPackage) error {
	const sql = `
		UPDATE order.cargo_packages
		SET name = @name,
			archived_at = @archived_at,
			updated_at = now()
		WHERE id = @id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":          pack.ID,
		"name":        pack.Name,
		"archived_at": pack.ArchivedAt,
	})
}

func (r *Repository) GetOffering(ctx context.Context, offeringID uuid.UUID) (model.Offering, error) {
	const sql = `
		SELECT *
		FROM order.offerings
		WHERE id = @id`

	rec, err := postgres.Get[record.Offering](ctx, r.pool, sql, pgx.NamedArgs{"id": offeringID}, postgres.WithNotFound(model.ErrOfferingNotFound))
	if err != nil {
		return model.Offering{}, err
	}

	return record.ToOffering(rec), nil
}

func (r *Repository) ListOfferings(ctx context.Context) ([]model.Offering, error) {
	const sql = `
		SELECT *
		FROM order.offerings`

	rec, err := postgres.Select[record.Offering](ctx, r.pool, sql, nil)
	if err != nil {
		return nil, err
	}

	return record.ToListOfferings(rec), nil
}

func (r *Repository) CreateOffering(ctx context.Context, offering model.Offering) error {
	const sql = `
		INSERT INTO order.offerings (id, name, price, modifiers)
		VALUES (@id, @name, @price, @modifiers)`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":        offering.ID,
		"name":      offering.Name,
		"price":     offering.Price,
		"modifiers": offering.Modifiers,
	})
}

func (r *Repository) UpdateOffering(ctx context.Context, offering model.Offering) error {
	const sql = `
		UPDATE order.offerings
		SET name = @name,
			price = @price,
			modifiers = @modifiers,
			archived_at = @archived_at,
			updated_at = now()
		WHERE id = @id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":          offering.ID,
		"name":        offering.Name,
		"price":       offering.Price,
		"modifiers":   offering.Modifiers,
		"archived_at": offering.ArchivedAt,
	})
}

func (r *Repository) GetPackingMaterial(ctx context.Context, materialID uuid.UUID) (model.PackingMaterial, error) {
	const sql = `
		SELECT *
		FROM order.packing_materials
		WHERE id = @id`

	rec, err := postgres.Get[record.PackingMaterial](ctx, r.pool, sql, pgx.NamedArgs{"id": materialID}, postgres.WithNotFound(model.ErrPackingMaterialNotFound))
	if err != nil {
		return model.PackingMaterial{}, err
	}

	return record.ToPackingMaterial(rec), nil
}

func (r *Repository) ListPackingMaterials(ctx context.Context) ([]model.PackingMaterial, error) {
	const sql = `
		SELECT *
		FROM order.packing_materials`

	rec, err := postgres.Select[record.PackingMaterial](ctx, r.pool, sql, nil)
	if err != nil {
		return nil, err
	}

	return record.ToListPackingMaterials(rec), nil
}

func (r *Repository) CreatePackingMaterial(ctx context.Context, material model.PackingMaterial) error {
	const sql = `
		INSERT INTO order.packing_materials (id, name, price)
		VALUES (@id, @name, @price)`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":    material.ID,
		"name":  material.Name,
		"price": material.Price,
	})
}

func (r *Repository) UpdatePackingMaterial(ctx context.Context, material model.PackingMaterial) error {
	const sql = `
		UPDATE order.packing_materials
		SET name = @name,
			price = @price,
			archived_at = @archived_at,
			updated_at = now()
		WHERE id = @id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":          material.ID,
		"name":        material.Name,
		"price":       material.Price,
		"archived_at": material.ArchivedAt,
	})
}

func (r *Repository) GetVehicle(ctx context.Context, vehicleID uuid.UUID) (model.Vehicle, error) {
	const sql = `
		SELECT *
		FROM order.vehicles
		WHERE vehicle_id = @vehicle_id`

	rec, err := postgres.Get[record.Vehicle](ctx, r.pool, sql, pgx.NamedArgs{"vehicle_id": vehicleID}, postgres.WithNotFound(model.ErrVehicleNotFound))
	if err != nil {
		return model.Vehicle{}, err
	}

	return record.ToVehicle(rec), nil
}

func (r *Repository) ListVehicles(ctx context.Context) ([]model.Vehicle, error) {
	const sql = `
		SELECT *
		FROM order.vehicles`

	rec, err := postgres.Select[record.Vehicle](ctx, r.pool, sql, nil)
	if err != nil {
		return nil, err
	}

	return record.ToListVehicles(rec), nil
}

func (r *Repository) CreateVehicle(ctx context.Context, vehicle model.Vehicle) error {
	const sql = `
		INSERT INTO order.vehicles (id, type, registration_number, length_meters, width_meters, height_meters, capacity_tonnes)
		VALUES (@id, @type, @registration_number, @length_meters, @width_meters, @height_meters, @capacity_tonnes)`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":                  vehicle.ID,
		"type":                vehicle.Type,
		"registration_number": vehicle.RegistrationNumber,
		"length_meters":       vehicle.LengthMeters,
		"width_meters":        vehicle.WidthMeters,
		"height_meters":       vehicle.HeightMeters,
		"capacity_tonnes":     vehicle.CapacityTonnes,
	}, postgres.WithExists(model.ErrVehicleExists))
}

func (r *Repository) UpdateVehicle(ctx context.Context, vehicle model.Vehicle) error {
	const sql = `
		UPDATE order.vehicles
		SET type = @type, 
			registration_number = @registration_number,
			length_meters = @length_meters,
			width_meters = @width_meters,
			height_meters = @height_meters,
			capacity_tonnes = @capacity_tonnes,
			archived_at = @archived_at,
			updated_at = now()
		WHERE id = @id`

	return postgres.Exec(ctx, r.pool, sql, pgx.NamedArgs{
		"id":                  vehicle.ID,
		"type":                vehicle.Type,
		"registration_number": vehicle.RegistrationNumber,
		"length_meters":       vehicle.LengthMeters,
		"width_meters":        vehicle.WidthMeters,
		"height_meters":       vehicle.HeightMeters,
		"capacity_tonnes":     vehicle.CapacityTonnes,
		"archived_at":         vehicle.ArchivedAt,
	}, postgres.WithExists(model.ErrVehicleExists))
}
