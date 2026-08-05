package model

import "errors"

var (
	ErrCargoPackageNotFound    = errors.New("cargo package not found")
	ErrOfferingNotFound        = errors.New("offering not found")
	ErrPackingMaterialNotFound = errors.New("packing material not found")
	ErrVehicleNotFound         = errors.New("vehicle not found")
	ErrVehicleExists           = errors.New("vehicle already exists")
)
