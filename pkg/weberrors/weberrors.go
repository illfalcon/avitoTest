package weberrors

import "github.com/pkg/errors"

var (
	ErrNotFound     = errors.New("requested item is not found")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrConflict     = errors.New("conflict with the current state of the resource")
	ErrBadRequest   = errors.New("malformed request")
	ErrForbidden    = errors.New("you are not allowed to perform this action")
)
