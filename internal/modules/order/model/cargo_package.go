package model

import (
	"errors"
	"time"

	"uuid"

	"github.com/samber/lo"
)

// CargoPackage Комплектация, то есть в чем транспортируются вещи, например: палет, россыпь или ящик.
type CargoPackage struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}

func NewCargoPackage(name string) (CargoPackage, error) {
	if name == "" {
		return CargoPackage{}, errors.New("name is required")
	}

	return CargoPackage{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (c *CargoPackage) Update(name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	c.Name = name

	return nil
}

func (c *CargoPackage) Archive() {
	if c.ArchivedAt == nil {
		c.ArchivedAt = lo.ToPtr(time.Now().UTC())
	}
}

func (c *CargoPackage) Restore() {
	c.ArchivedAt = nil
}
