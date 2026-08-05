package model

import (
	"errors"
	"time"

	"uuid"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

// Offering Дополнительная услуга к заказу
type Offering struct {
	ID         uuid.UUID          `json:"id"`
	Name       string             `json:"name"`
	Price      decimal.Decimal    `json:"price"`
	Modifiers  []OfferingModifier `json:"modifiers"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	ArchivedAt *time.Time         `json:"archived_at"`
}

func NewOffering(name string, price decimal.Decimal, modifiers []OfferingModifier) (Offering, error) {
	if name == "" {
		return Offering{}, errors.New("name is required")
	}
	if price.IsNegative() {
		return Offering{}, errors.New("price is negative")
	}

	return Offering{
		ID:        uuid.New(),
		Name:      name,
		Price:     price,
		Modifiers: modifiers,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (o *Offering) Update(name string, price decimal.Decimal, modifiers []OfferingModifier) error {
	if name == "" {
		return errors.New("name is required")
	}
	if price.IsNegative() {
		return errors.New("price is negative")
	}

	o.Name = name
	o.Price = price
	o.Modifiers = modifiers

	return nil
}

func (o *Offering) Archive() {
	if o.ArchivedAt == nil {
		o.ArchivedAt = lo.ToPtr(time.Now().UTC())
	}
}

func (o *Offering) Restore() {
	o.ArchivedAt = nil
}
