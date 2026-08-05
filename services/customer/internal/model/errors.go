package model

import "errors"

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCustomerExists   = errors.New("customer already exists")
)
