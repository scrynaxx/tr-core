package record

import "github.com/scrynaxx/tr-core/internal/modules/order/model"

func ToCargoPackage(rec CargoPackage) model.CargoPackage {
	return model.CargoPackage{
		ID:         rec.ID,
		Name:       rec.Name,
		CreatedAt:  rec.CreatedAt,
		UpdatedAt:  rec.UpdatedAt,
		ArchivedAt: rec.ArchivedAt,
	}
}

func ToListCargoPackages(recs []CargoPackage) []model.CargoPackage {
	list := make([]model.CargoPackage, len(recs))
	for i, rec := range recs {
		list[i] = ToCargoPackage(rec)
	}

	return list
}

func ToOffering(rec Offering) model.Offering {
	return model.Offering{
		ID:         rec.ID,
		Name:       rec.Name,
		Price:      rec.Price,
		Modifiers:  rec.Modifiers,
		CreatedAt:  rec.CreatedAt,
		UpdatedAt:  rec.UpdatedAt,
		ArchivedAt: rec.ArchivedAt,
	}
}

func ToListOfferings(recs []Offering) []model.Offering {
	list := make([]model.Offering, len(recs))
	for i, rec := range recs {
		list[i] = ToOffering(rec)
	}

	return list
}

func ToPackingMaterial(rec PackingMaterial) model.PackingMaterial {
	return model.PackingMaterial{
		ID:         rec.ID,
		Name:       rec.Name,
		Price:      rec.Price,
		CreatedAt:  rec.CreatedAt,
		UpdatedAt:  rec.UpdatedAt,
		ArchivedAt: rec.ArchivedAt,
	}
}

func ToListPackingMaterials(recs []PackingMaterial) []model.PackingMaterial {
	list := make([]model.PackingMaterial, len(recs))
	for i, rec := range recs {
		list[i] = ToPackingMaterial(rec)
	}

	return list
}

func ToVehicle(rec Vehicle) model.Vehicle {
	return model.Vehicle{
		ID:                 rec.ID,
		Type:               rec.Type,
		RegistrationNumber: rec.RegistrationNumber,
		LengthMeters:       rec.LengthMeters,
		WidthMeters:        rec.WidthMeters,
		HeightMeters:       rec.HeightMeters,
		CapacityTonnes:     rec.CapacityTonnes,
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          rec.UpdatedAt,
		ArchivedAt:         rec.ArchivedAt,
	}
}

func ToListVehicles(recs []Vehicle) []model.Vehicle {
	list := make([]model.Vehicle, len(recs))
	for i, rec := range recs {
		list[i] = ToVehicle(rec)
	}

	return list
}
