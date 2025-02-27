package models

type Inventory struct {
	ID         uint32   `json:"id"`
	Name       string   `json:"name"`
	Quantity   int      `json:"quantity"`
	Unit       string   `json:"unit"`
	Categories []string `json:"categories"`
}

func NewInventoryValidator(inventory *Inventory) inventoryValidator {
	return inventoryValidator{
		make(map[string]string),
		inventory,
	}
}

type inventoryValidator struct {
	validator map[string]string
	inventory *Inventory
}

func (v *inventoryValidator) Validate() (map[string]string, bool) {
	if v.inventory.Name == "" {
		v.validator["Name"] = "missing Name"
	}
	if v.inventory.Unit == "" {
		v.validator["Unit"] = "missing Unit"
	}

	if len(v.validator) > 0 {
		return v.validator, false
	} else {
		return nil, true
	}
}
