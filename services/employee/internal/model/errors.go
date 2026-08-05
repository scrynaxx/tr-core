package model

import "errors"

var (
	ErrEmployeeNotFound    = errors.New("employee not found")
	ErrEmployeeExists      = errors.New("employee already exists")
	ErrEmailExists         = errors.New("email already exists")
	ErrPassportNotFound    = errors.New("passport not found")
	ErrCredentialsNotFound = errors.New("credentials not found")
)
