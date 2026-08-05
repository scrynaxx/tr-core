package model

// VehicleType Тип транспорта.
type VehicleType string

const (
	// VehicleTypeAwning Покрытый тентом.
	VehicleTypeAwning VehicleType = "awning"

	// VehicleTypeFridge Холодильник (рефрижератор).
	VehicleTypeFridge VehicleType = "fridge"

	// VehicleTypeThermos Термос.
	VehicleTypeThermos VehicleType = "thermos"
)
