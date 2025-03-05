package models

type MenuItem struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Price       float64             `json:"price"`
	Inventory   []MenuItemInventory `json:"inventory"`
}

type MenuItemInventory struct {
	InventoryID int     `json:"inventory_id"`
	Quantity    float64 `json:"quantity"`
}
