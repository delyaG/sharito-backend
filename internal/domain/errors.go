package domain

import "fmt"

var (
	// StatusBadRequest
	ErrInvalidInputData = fmt.Errorf("invalid input data")
	ErrNoSuchUser       = fmt.Errorf("no such user error")

	// StatusInternalServerError
	ErrInternalSecurity = fmt.Errorf("internal security error")
	ErrInternalDatabase = fmt.Errorf("internal database error")
	ErrJWT              = fmt.Errorf("jwt creating error")

	// StatusUnauthorized
	ErrUnauthorized = fmt.Errorf("unauthorized")
)
