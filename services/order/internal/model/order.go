package model

import (
	"time"

	"github.com/google/uuid"
)

// Order Заказ на транспортировку.
type Order struct {
	ID              uuid.UUID      `json:"id"`
	CustomerID      uuid.UUID      `json:"customer_id"`
	ManagerID       uuid.UUID      `json:"manager_id"`
	VehicleID       uuid.UUID      `json:"vehicle_id"`
	Type            OrderType      `json:"type"`
	Status          OrderStatus    `json:"status"`
	LoadAddresses   []string       `json:"load_addresses"`
	UnloadAddresses []string       `json:"unload_addresses"`
	Payment         OrderPayment   `json:"payment"`
	Comments        []OrderComment `json:"comments"`
	OvertimeHours   int            `json:"overtime_hours"`

	DetailsOffice  *OrderOffice  `json:"office,omitempty"`
	DetailsFreight *OrderFreight `json:"freight,omitempty"`

	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderOffice struct {
	Offerings []OrderOffering `json:"offers"`
	Packings  []OrderPacking  `json:"packings"`
}

type OrderFreight struct {
	SeatsCount        int                 `json:"seats_count"`
	TotalVolume       float64             `json:"total_volume"`
	TotalWeight       float64             `json:"total_weight"`
	LoadType          LoadType            `json:"load_type"`
	UnloadType        LoadType            `json:"unload_type"`
	CargoPackages     []OrderCargoPackage `json:"cargo_packages"`
	SpecialConditions *string             `json:"special_conditions"`
}
