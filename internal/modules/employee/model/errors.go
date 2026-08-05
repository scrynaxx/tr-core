package model

import (
	"errors"
)

var (
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrPassportNotFound = errors.New("passport not found")
)
