package dictionary

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/order/model"
	"github.com/scrynaxx/tr-core/internal/modules/order/repository"
	"github.com/scrynaxx/tr-core/internal/modules/order/usecase"
	"github.com/shopspring/decimal"
)

type UseCase struct {
	dictionaryRepo repository.Dictionary
}

func New(dictionaryRepo repository.Dictionary) usecase.Dictionary {
	return &UseCase{
		dictionaryRepo: dictionaryRepo,
	}
}

func (uc *UseCase) GetCargoPackage(ctx context.Context, packageID uuid.UUID) (model.CargoPackage, error) {
	return uc.dictionaryRepo.GetCargoPackage(ctx, packageID)
}

func (uc *UseCase) ListCargoPackages(ctx context.Context) ([]model.CargoPackage, error) {
	return uc.dictionaryRepo.ListCargoPackages(ctx)
}

func (uc *UseCase) CreateCargoPackage(ctx context.Context, name string) error {
	cargo, err := model.NewCargoPackage(name)
	if err != nil {
		return err
	}

	return uc.dictionaryRepo.CreateCargoPackage(ctx, cargo)
}

func (uc *UseCase) UpdateCargoPackage(ctx context.Context, packageID uuid.UUID, name string) error {
	cargo, err := uc.dictionaryRepo.GetCargoPackage(ctx, packageID)
	if err != nil {
		return err
	}

	if err = cargo.Update(name); err != nil {
		return err
	}

	return uc.dictionaryRepo.UpdateCargoPackage(ctx, cargo)
}

func (uc *UseCase) ArchiveCargoPackage(ctx context.Context, packageID uuid.UUID) error {
	cargo, err := uc.dictionaryRepo.GetCargoPackage(ctx, packageID)
	if err != nil {
		return err
	}

	cargo.Archive()

	return uc.dictionaryRepo.UpdateCargoPackage(ctx, cargo)
}

func (uc *UseCase) RestoreCargoPackage(ctx context.Context, packageID uuid.UUID) error {
	cargo, err := uc.dictionaryRepo.GetCargoPackage(ctx, packageID)
	if err != nil {
		return err
	}

	cargo.Restore()

	return uc.dictionaryRepo.UpdateCargoPackage(ctx, cargo)
}

func (uc *UseCase) GetOffering(ctx context.Context, offeringID uuid.UUID) (model.Offering, error) {
	return uc.dictionaryRepo.GetOffering(ctx, offeringID)
}

func (uc *UseCase) ListOfferings(ctx context.Context) ([]model.Offering, error) {
	return uc.dictionaryRepo.ListOfferings(ctx)
}

func (uc *UseCase) CreateOffering(ctx context.Context, name string, price decimal.Decimal, modifiers []model.OfferingModifier) error {
	offering, err := model.NewOffering(name, price, modifiers)
	if err != nil {
		return err
	}

	return uc.dictionaryRepo.CreateOffering(ctx, offering)
}

func (uc *UseCase) UpdateOffering(ctx context.Context, offeringID uuid.UUID, name string, price decimal.Decimal, modifiers []model.OfferingModifier) error {
	offering, err := uc.dictionaryRepo.GetOffering(ctx, offeringID)
	if err != nil {
		return err
	}

	if err = offering.Update(name, price, modifiers); err != nil {
		return err
	}

	return uc.dictionaryRepo.UpdateOffering(ctx, offering)
}

func (uc *UseCase) ArchiveOffering(ctx context.Context, offeringID uuid.UUID) error {
	offering, err := uc.dictionaryRepo.GetOffering(ctx, offeringID)
	if err != nil {
		return err
	}

	offering.Archive()

	return uc.dictionaryRepo.UpdateOffering(ctx, offering)
}

func (uc *UseCase) RestoreOffering(ctx context.Context, offeringID uuid.UUID) error {
	offering, err := uc.dictionaryRepo.GetOffering(ctx, offeringID)
	if err != nil {
		return err
	}

	offering.Restore()

	return uc.dictionaryRepo.UpdateOffering(ctx, offering)
}

func (uc *UseCase) GetPackingMaterial(ctx context.Context, materialID uuid.UUID) (model.PackingMaterial, error) {
	return uc.dictionaryRepo.GetPackingMaterial(ctx, materialID)
}

func (uc *UseCase) ListPackingMaterials(ctx context.Context) ([]model.PackingMaterial, error) {
	return uc.dictionaryRepo.ListPackingMaterials(ctx)
}

func (uc *UseCase) CreatePackingMaterial(ctx context.Context, name string, price decimal.Decimal) error {
	packing, err := model.NewPackingMaterial(name, price)
	if err != nil {
		return err
	}

	return uc.dictionaryRepo.CreatePackingMaterial(ctx, packing)
}

func (uc *UseCase) UpdatePackingMaterial(ctx context.Context, materialID uuid.UUID, name string, price decimal.Decimal) error {
	packing, err := uc.dictionaryRepo.GetPackingMaterial(ctx, materialID)
	if err != nil {
		return err
	}

	if err = packing.Update(name, price); err != nil {
		return err
	}

	return uc.dictionaryRepo.UpdatePackingMaterial(ctx, packing)
}

func (uc *UseCase) ArchivePackingMaterial(ctx context.Context, materialID uuid.UUID) error {
	packing, err := uc.dictionaryRepo.GetPackingMaterial(ctx, materialID)
	if err != nil {
		return err
	}

	packing.Archive()

	return uc.dictionaryRepo.UpdatePackingMaterial(ctx, packing)
}

func (uc *UseCase) RestorePackingMaterial(ctx context.Context, materialID uuid.UUID) error {
	packing, err := uc.dictionaryRepo.GetPackingMaterial(ctx, materialID)
	if err != nil {
		return err
	}

	packing.Restore()

	return uc.dictionaryRepo.UpdatePackingMaterial(ctx, packing)
}

func (uc *UseCase) GetVehicle(ctx context.Context, vehicleID uuid.UUID) (model.Vehicle, error) {
	return uc.dictionaryRepo.GetVehicle(ctx, vehicleID)
}

func (uc *UseCase) ListVehicles(ctx context.Context) ([]model.Vehicle, error) {
	return uc.dictionaryRepo.ListVehicles(ctx)
}

func (uc *UseCase) CreateVehicle(ctx context.Context, input model.VehicleInput) error {
	vehicle, err := model.NewVehicle(input)
	if err != nil {
		return err
	}

	return uc.dictionaryRepo.CreateVehicle(ctx, vehicle)
}

func (uc *UseCase) UpdateVehicle(ctx context.Context, vehicleID uuid.UUID, input model.VehicleInput) error {
	vehicle, err := uc.dictionaryRepo.GetVehicle(ctx, vehicleID)
	if err != nil {
		return err
	}

	if err = vehicle.Update(input); err != nil {
		return err
	}

	return uc.dictionaryRepo.UpdateVehicle(ctx, vehicle)
}

func (uc *UseCase) ArchiveVehicle(ctx context.Context, vehicleID uuid.UUID) error {
	vehicle, err := uc.dictionaryRepo.GetVehicle(ctx, vehicleID)
	if err != nil {
		return err
	}

	vehicle.Archive()

	return uc.dictionaryRepo.UpdateVehicle(ctx, vehicle)
}

func (uc *UseCase) RestoreVehicle(ctx context.Context, vehicleID uuid.UUID) error {
	vehicle, err := uc.dictionaryRepo.GetVehicle(ctx, vehicleID)
	if err != nil {
		return err
	}

	vehicle.Restore()

	return uc.dictionaryRepo.UpdateVehicle(ctx, vehicle)
}
