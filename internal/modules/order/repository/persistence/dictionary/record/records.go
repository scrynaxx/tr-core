package record

import (
	"time"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/order/model"
	"github.com/shopspring/decimal"
)

type CargoPackage struct {
	ID         uuid.UUID  `db:"id"`
	Name       string     `db:"name"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	ArchivedAt *time.Time `db:"archived_at"`
}

type Offering struct {
	ID         uuid.UUID                `db:"id"`
	Name       string                   `db:"name"`
	Price      decimal.Decimal          `db:"price"`
	Modifiers  []model.OfferingModifier `db:"modifiers"`
	CreatedAt  time.Time                `db:"created_at"`
	UpdatedAt  time.Time                `db:"updated_at"`
	ArchivedAt *time.Time               `db:"archived_at"`
}

type PackingMaterial struct {
	ID         uuid.UUID       `db:"id"`
	Name       string          `db:"name"`
	Price      decimal.Decimal `db:"price"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
	ArchivedAt *time.Time      `db:"archived_at"`
}

type Vehicle struct {
	ID                 uuid.UUID         `db:"id"`
	Type               model.VehicleType `db:"type"`
	RegistrationNumber *string           `db:"registration_number"`
	LengthMeters       float64           `db:"length_meters"`
	WidthMeters        float64           `db:"width_meters"`
	HeightMeters       float64           `db:"height_meters"`
	CapacityTonnes     float64           `db:"capacity_tonnes"`
	CreatedAt          time.Time         `db:"created_at"`
	UpdatedAt          time.Time         `db:"updated_at"`
	ArchivedAt         *time.Time        `db:"archived_at"`
}
