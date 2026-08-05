package model

import (
	"time"

	"github.com/google/uuid"
)

// CargoPackage Комплектация, то есть в чем транспортируются вещи, например: палет, россыпь или ящик.
type CargoPackage struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}
