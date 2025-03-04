package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no record")

	//Inventory errors
	ErrDuplicateInventory = errors.New("models: duplicate inventory")
	ErrNegativeQuantity   = errors.New("models: positive_quantity constraint violation")
	ErrMissingFields      = errors.New("models: missing fields")
	ErrInvalidEnumType    = errors.New("models: invalid enum type")
	ErrInvalidID          = errors.New("models: id is not valid int")
)
