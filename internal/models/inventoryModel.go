package models

type Inventory struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Quantity   int      `json:"quantity"`
	Unit       string   `json:"unit"`
	Categories []string `json:"categories"`
}

type inventoryValidator struct {
	validator map[string]string
	inventory *Inventory
}

func NewInventoryValidator(inventory *Inventory) *inventoryValidator {
	return &inventoryValidator{
		make(map[string]string),
		inventory,
	}
}

func (v *inventoryValidator) Validate() map[string]string {
	if v.inventory.Name == "" {
		v.validator["Name"] = "missing Name"
	}
	if v.inventory.Unit == "" {
		v.validator["Unit"] = "missing Unit"
	}

	if len(v.validator) > 0 {
		return v.validator
	} else {
		return nil
	}
}
