package contract

import (
	"time"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/order/model"
	"github.com/shopspring/decimal"
)

type CreateCargoPackageRequest struct {
	Name string `json:"name" validate:"required"`
}

type GetCargoPackageRequest struct {
	PackageID uuid.UUID `param:"package_id" validate:"required,uuid"`
}

type UpdateCargoPackageRequest struct {
	PackageID uuid.UUID `param:"package_id" validate:"required,uuid"`
	Name      string    `json:"name" validate:"required"`
}

type ArchiveCargoPackageRequest struct {
	PackageID uuid.UUID `param:"package_id" validate:"required,uuid"`
}

type RestoreCargoPackageRequest struct {
	PackageID uuid.UUID `param:"package_id" validate:"required,uuid"`
}

type ListCargoPackagesResponse struct {
	Packages []CargoPackage `json:"packages"`
}

type CreateOfferingRequest struct {
	Name      string                   `json:"name" validate:"required"`
	Price     decimal.Decimal          `json:"price" validate:"required"`
	Modifiers []model.OfferingModifier `json:"modifiers"`
}

type GetOfferingRequest struct {
	OfferingID uuid.UUID `param:"offering_id" validate:"required,uuid"`
}

type UpdateOfferingRequest struct {
	OfferingID uuid.UUID                `param:"offering_id" validate:"required,uuid"`
	Name       string                   `json:"name" validate:"required"`
	Price      decimal.Decimal          `json:"price" validate:"required"`
	Modifiers  []model.OfferingModifier `json:"modifiers"`
}

type ArchiveOfferingRequest struct {
	OfferingID uuid.UUID `param:"offering_id" validate:"required,uuid"`
}

type RestoreOfferingRequest struct {
	OfferingID uuid.UUID `param:"offering_id" validate:"required,uuid"`
}

type ListOfferingsResponse struct {
	Offerings []Offering `json:"offerings"`
}

type CreatePackingMaterialRequest struct {
	Name  string          `json:"name" validate:"required"`
	Price decimal.Decimal `json:"price" validate:"required"`
}

type GetPackingMaterialRequest struct {
	MaterialID uuid.UUID `param:"material_id" validate:"required,uuid"`
}

type UpdatePackingMaterialRequest struct {
	MaterialID uuid.UUID       `param:"material_id" validate:"required,uuid"`
	Name       string          `json:"name" validate:"required"`
	Price      decimal.Decimal `json:"price" validate:"required"`
}

type ArchivePackingMaterialRequest struct {
	MaterialID uuid.UUID `param:"material_id" validate:"required,uuid"`
}

type RestorePackingMaterialRequest struct {
	MaterialID uuid.UUID `param:"material_id" validate:"required,uuid"`
}

type ListPackingMaterialsResponse struct {
	Materials []PackingMaterial `json:"materials"`
}

type CreateVehicleRequest struct {
	Type               model.VehicleType `json:"type" validate:"required"`
	RegistrationNumber *string           `json:"registration_number"`
	LengthMeters       float64           `json:"length_meters" validate:"required"`
	WidthMeters        float64           `json:"width_meters" validate:"required"`
	HeightMeters       float64           `json:"height_meters" validate:"required"`
	CapacityTonnes     float64           `json:"capacity_tonnes" validate:"required"`
}

type GetVehicleRequest struct {
	VehicleID uuid.UUID `param:"vehicle_id" validate:"required,uuid"`
}

type UpdateVehicleRequest struct {
	VehicleID          uuid.UUID         `param:"vehicle_id" validate:"required,uuid"`
	Type               model.VehicleType `json:"type" validate:"required"`
	RegistrationNumber *string           `json:"registration_number"`
	LengthMeters       float64           `json:"length_meters" validate:"required"`
	WidthMeters        float64           `json:"width_meters" validate:"required"`
	HeightMeters       float64           `json:"height_meters" validate:"required"`
	CapacityTonnes     float64           `json:"capacity_tonnes" validate:"required"`
}

type ArchiveVehicleRequest struct {
	VehicleID uuid.UUID `param:"vehicle_id" validate:"required,uuid"`
}

type RestoreVehicleRequest struct {
	VehicleID uuid.UUID `param:"vehicle_id" validate:"required,uuid"`
}

type ListVehiclesResponse struct {
	Vehicles []Vehicle `json:"vehicles"`
}

type CargoPackage struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}

type Offering struct {
	ID         uuid.UUID                `json:"id"`
	Name       string                   `json:"name"`
	Price      decimal.Decimal          `json:"price"`
	Modifiers  []model.OfferingModifier `json:"modifiers"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
	ArchivedAt *time.Time               `json:"archived_at"`
}

type PackingMaterial struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Price      decimal.Decimal `json:"price"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	ArchivedAt *time.Time      `json:"archived_at"`
}

type Vehicle struct {
	ID                 uuid.UUID         `json:"id"`
	Type               model.VehicleType `json:"type"`
	RegistrationNumber *string           `json:"registration_number"`
	LengthMeters       float64           `json:"length_meters"`
	WidthMeters        float64           `json:"width_meters"`
	HeightMeters       float64           `json:"height_meters"`
	CapacityTonnes     float64           `json:"capacity_tonnes"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
	ArchivedAt         *time.Time        `json:"archived_at"`
}
