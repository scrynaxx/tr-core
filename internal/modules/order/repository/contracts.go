package repository

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/order/model"
)

type Dictionary interface {
	GetCargoPackage(ctx context.Context, packageID uuid.UUID) (model.CargoPackage, error)
	ListCargoPackages(ctx context.Context) ([]model.CargoPackage, error)
	CreateCargoPackage(ctx context.Context, pack model.CargoPackage) error
	UpdateCargoPackage(ctx context.Context, pack model.CargoPackage) error

	GetOffering(ctx context.Context, offeringID uuid.UUID) (model.Offering, error)
	ListOfferings(ctx context.Context) ([]model.Offering, error)
	CreateOffering(ctx context.Context, offering model.Offering) error
	UpdateOffering(ctx context.Context, offering model.Offering) error

	GetPackingMaterial(ctx context.Context, materialID uuid.UUID) (model.PackingMaterial, error)
	ListPackingMaterials(ctx context.Context) ([]model.PackingMaterial, error)
	CreatePackingMaterial(ctx context.Context, material model.PackingMaterial) error
	UpdatePackingMaterial(ctx context.Context, material model.PackingMaterial) error

	GetVehicle(ctx context.Context, vehicleID uuid.UUID) (model.Vehicle, error)
	ListVehicles(ctx context.Context) ([]model.Vehicle, error)
	CreateVehicle(ctx context.Context, vehicle model.Vehicle) error
	UpdateVehicle(ctx context.Context, vehicle model.Vehicle) error
}
