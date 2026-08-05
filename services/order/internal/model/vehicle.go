package model

import (
	"time"

	"github.com/google/uuid"
)

// Vehicle Транспорт для транспортировки.
type Vehicle struct {
	ID                 uuid.UUID   `json:"id"`
	Type               VehicleType `json:"type"`
	RegistrationNumber string      `json:"registration_number"`
	LengthMeters       float64     `json:"length_meters"`
	WidthMeters        float64     `json:"width_meters"`
	HeightMeters       float64     `json:"height_meters"`
	CapacityTonnes     float64     `json:"capacity_tonnes"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	ArchivedAt         *time.Time  `json:"archived_at"`
}
