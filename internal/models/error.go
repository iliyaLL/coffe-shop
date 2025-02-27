package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no record")
	ErrDuplicateInventory = errors.New("models: duplicate inventory")
	ErrNegativeQuantity   = errors.New("models: positive_quantity constraint violation")
	ErrMissingFields      = errors.New("models: missing fields")
	ErrInvalidEnumType    = errors.New("models: invalid enum type")
)
