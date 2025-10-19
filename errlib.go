package ecstore

import (
	"errors"
)

var (
	ErrNoEntitiesProvided   = errors.New("store: no entities provided to the Add function")
	ErrInvalidEntityPointer = errors.New("store: entity must be a non-nil pointer to a struct")
)