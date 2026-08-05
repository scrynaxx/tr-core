package model

import (
	"errors"
	"time"

	"uuid"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

// PackingMaterial Упаковка / упаковочный материал.
type PackingMaterial struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Price      decimal.Decimal `json:"price"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	ArchivedAt *time.Time      `json:"archived_at"`
}

func NewPackingMaterial(name string, price decimal.Decimal) (PackingMaterial, error) {
	if name == "" {
		return PackingMaterial{}, errors.New("name is required")
	}
	if price.IsNegative() {
		return PackingMaterial{}, errors.New("price is negative")
	}

	return PackingMaterial{
		ID:        uuid.New(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (p *PackingMaterial) Update(name string, price decimal.Decimal) error {
	if name == "" {
		return errors.New("name is required")
	}
	if price.IsNegative() {
		return errors.New("price is negative")
	}

	p.Name = name
	p.Price = price

	return nil
}

func (p *PackingMaterial) Archive() {
	if p.ArchivedAt == nil {
		p.ArchivedAt = lo.ToPtr(time.Now().UTC())
	}
}

func (p *PackingMaterial) Restore() {
	p.ArchivedAt = nil
}
