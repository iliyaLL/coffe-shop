package models

import "errors"

var (
	ErrDuplicateInventory = errors.New("models: duplicate inventory")
	ErrNoRecord           = errors.New("models: no record")
	ErrNegativeQuantity   = errors.New("models: positive_quantity constraint violation")
	ErrMissingFields      = errors.New("models: missing fields")
)
