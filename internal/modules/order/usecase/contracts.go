package usecase

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/order/model"
	"github.com/shopspring/decimal"
)

type Dictionary interface {
	GetCargoPackage(ctx context.Context, packageID uuid.UUID) (model.CargoPackage, error)
	ListCargoPackages(ctx context.Context) ([]model.CargoPackage, error)
	CreateCargoPackage(ctx context.Context, name string) error
	UpdateCargoPackage(ctx context.Context, packageID uuid.UUID, name string) error
	ArchiveCargoPackage(ctx context.Context, packageID uuid.UUID) error
	RestoreCargoPackage(ctx context.Context, packageID uuid.UUID) error

	GetOffering(ctx context.Context, offeringID uuid.UUID) (model.Offering, error)
	ListOfferings(ctx context.Context) ([]model.Offering, error)
	CreateOffering(ctx context.Context, name string, price decimal.Decimal, modifiers []model.OfferingModifier) error
	UpdateOffering(ctx context.Context, offeringID uuid.UUID, name string, price decimal.Decimal, modifiers []model.OfferingModifier) error
	ArchiveOffering(ctx context.Context, offeringID uuid.UUID) error
	RestoreOffering(ctx context.Context, offeringID uuid.UUID) error

	GetPackingMaterial(ctx context.Context, materialID uuid.UUID) (model.PackingMaterial, error)
	ListPackingMaterials(ctx context.Context) ([]model.PackingMaterial, error)
	CreatePackingMaterial(ctx context.Context, name string, price decimal.Decimal) error
	UpdatePackingMaterial(ctx context.Context, materialID uuid.UUID, name string, price decimal.Decimal) error
	ArchivePackingMaterial(ctx context.Context, materialID uuid.UUID) error
	RestorePackingMaterial(ctx context.Context, materialID uuid.UUID) error

	GetVehicle(ctx context.Context, vehicleID uuid.UUID) (model.Vehicle, error)
	ListVehicles(ctx context.Context) ([]model.Vehicle, error)
	CreateVehicle(ctx context.Context, input model.VehicleInput) error
	UpdateVehicle(ctx context.Context, vehicleID uuid.UUID, input model.VehicleInput) error
	ArchiveVehicle(ctx context.Context, vehicleID uuid.UUID) error
	RestoreVehicle(ctx context.Context, vehicleID uuid.UUID) error
}
