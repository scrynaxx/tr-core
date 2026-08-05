package model

import (
	"uuid"
)

// OrderCargoPackage Комплектация привязанная к заказу.
type OrderCargoPackage struct {
	CargoPackageID *uuid.UUID `json:"cargo_package_id"`
	Name           string     `json:"name"`
}
