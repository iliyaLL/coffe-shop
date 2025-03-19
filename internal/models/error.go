package models

import "errors"

var (
	ErrNoRecord         = errors.New("models: no record")
	ErrNegativeQuantity = errors.New("models: positive quantity constraint violation")
	ErrNegativePrice    = errors.New("models: positive price constraint violation")
	ErrMissingFields    = errors.New("models: missing fields")
	ErrInvalidID        = errors.New("models: id is not valid int")

	//Inventory errors
	ErrDuplicateInventory       = errors.New("models: duplicate inventory")
	ErrInvalidEnumTypeInventory = errors.New("models: invalid enum type. Supported types: shots, ml, g, units")

	//Menu errors
	ErrDuplicateMenuItem                 = errors.New("models: duplicate menu item")
	ErrForeignKeyConstraintMenuInventory = errors.New("models: inventory does not exist")
)
