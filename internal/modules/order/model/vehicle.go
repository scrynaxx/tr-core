package model

import (
	"errors"
	"time"

	"uuid"

	"github.com/samber/lo"
)

// Vehicle Транспорт для транспортировки.
type Vehicle struct {
	ID                 uuid.UUID   `json:"id"`
	Type               VehicleType `json:"type"`
	RegistrationNumber *string     `json:"registration_number"`
	LengthMeters       float64     `json:"length_meters"`
	WidthMeters        float64     `json:"width_meters"`
	HeightMeters       float64     `json:"height_meters"`
	CapacityTonnes     float64     `json:"capacity_tonnes"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	ArchivedAt         *time.Time  `json:"archived_at"`
}

type VehicleInput struct {
	Type               VehicleType
	RegistrationNumber *string
	LengthMeters       float64
	WidthMeters        float64
	HeightMeters       float64
	CapacityTonnes     float64
}

func NewVehicle(input VehicleInput) (Vehicle, error) {
	if input.Type == "" {
		return Vehicle{}, errors.New("type is required")
	}
	if input.LengthMeters == 0 {
		return Vehicle{}, errors.New("length meters is required")
	}
	if input.WidthMeters == 0 {
		return Vehicle{}, errors.New("width meters is required")
	}
	if input.HeightMeters == 0 {
		return Vehicle{}, errors.New("height meters is required")
	}
	if input.CapacityTonnes == 0 {
		return Vehicle{}, errors.New("capacity tonnes is required")
	}

	return Vehicle{
		ID:                 uuid.New(),
		Type:               input.Type,
		RegistrationNumber: input.RegistrationNumber,
		LengthMeters:       input.LengthMeters,
		WidthMeters:        input.WidthMeters,
		HeightMeters:       input.HeightMeters,
		CapacityTonnes:     input.CapacityTonnes,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
		ArchivedAt:         nil,
	}, nil
}

func (v *Vehicle) Update(input VehicleInput) error {
	if input.Type == "" {
		return errors.New("type is required")
	}
	if input.LengthMeters == 0 {
		return errors.New("length meters is required")
	}
	if input.WidthMeters == 0 {
		return errors.New("width meters is required")
	}
	if input.HeightMeters == 0 {
		return errors.New("height meters is required")
	}
	if input.CapacityTonnes == 0 {
		return errors.New("capacity tonnes is required")
	}

	v.Type = input.Type
	v.RegistrationNumber = input.RegistrationNumber
	v.LengthMeters = input.LengthMeters
	v.WidthMeters = input.WidthMeters
	v.HeightMeters = input.HeightMeters
	v.CapacityTonnes = input.CapacityTonnes

	return nil
}

func (v *Vehicle) Archive() {
	if v.ArchivedAt == nil {
		v.ArchivedAt = lo.ToPtr(time.Now().UTC())
	}
}

func (v *Vehicle) Restore() {
	v.ArchivedAt = nil
}
