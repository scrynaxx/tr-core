package contract

import "github.com/scrynaxx/tr-core/internal/modules/order/model"

func ToCargoPackage(pack model.CargoPackage) CargoPackage {
	return CargoPackage{
		ID:         pack.ID,
		Name:       pack.Name,
		CreatedAt:  pack.CreatedAt,
		UpdatedAt:  pack.UpdatedAt,
		ArchivedAt: pack.ArchivedAt,
	}
}

func ToListCargoPackagesResponse(packs []model.CargoPackage) ListCargoPackagesResponse {
	res := ListCargoPackagesResponse{
		Packages: make([]CargoPackage, len(packs)),
	}

	for i, pack := range packs {
		res.Packages[i] = ToCargoPackage(pack)
	}

	return res
}

func ToOffering(offering model.Offering) Offering {
	return Offering{
		ID:         offering.ID,
		Name:       offering.Name,
		Price:      offering.Price,
		Modifiers:  offering.Modifiers,
		CreatedAt:  offering.CreatedAt,
		UpdatedAt:  offering.UpdatedAt,
		ArchivedAt: offering.ArchivedAt,
	}
}

func ToListOfferingsResponse(offerings []model.Offering) ListOfferingsResponse {
	res := ListOfferingsResponse{
		Offerings: make([]Offering, len(offerings)),
	}

	for i, offering := range offerings {
		res.Offerings[i] = ToOffering(offering)
	}

	return res
}

func ToPackingMaterial(material model.PackingMaterial) PackingMaterial {
	return PackingMaterial{
		ID:         material.ID,
		Name:       material.Name,
		Price:      material.Price,
		CreatedAt:  material.CreatedAt,
		UpdatedAt:  material.UpdatedAt,
		ArchivedAt: material.ArchivedAt,
	}
}

func ToListPackingMaterialsResponse(materials []model.PackingMaterial) ListPackingMaterialsResponse {
	res := ListPackingMaterialsResponse{
		Materials: make([]PackingMaterial, len(materials)),
	}

	for i, material := range materials {
		res.Materials[i] = ToPackingMaterial(material)
	}

	return res
}

func ToVehicle(vehicle model.Vehicle) Vehicle {
	return Vehicle{
		ID:                 vehicle.ID,
		Type:               vehicle.Type,
		RegistrationNumber: vehicle.RegistrationNumber,
		LengthMeters:       vehicle.LengthMeters,
		WidthMeters:        vehicle.WidthMeters,
		HeightMeters:       vehicle.HeightMeters,
		CapacityTonnes:     vehicle.CapacityTonnes,
		CreatedAt:          vehicle.CreatedAt,
		UpdatedAt:          vehicle.UpdatedAt,
		ArchivedAt:         vehicle.ArchivedAt,
	}
}

func ToListVehiclesResponse(vehicles []model.Vehicle) ListVehiclesResponse {
	res := ListVehiclesResponse{
		Vehicles: make([]Vehicle, len(vehicles)),
	}

	for i, vehicle := range vehicles {
		res.Vehicles[i] = ToVehicle(vehicle)
	}

	return res
}
