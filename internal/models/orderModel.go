package models

import (
	"strconv"
	"time"
)

type Jsonb map[string]interface{}

type Order struct {
	ID                  int         `json:"id"`
	CustomerName        string      `json:"customer_name"`
	Status              string      `json:"status"`
	CreatedAt           time.Time   `json:"created_at"`
	CustomerPreferences Jsonb       `json:"customer_preferences"`
	Items               []OrderItem `json:"items"`
}

type OrderItem struct {
	MenuID   int `json:"menu_id"`
	Quantity int `json:"quantity"`
}

type orderValidator struct {
	errors map[string]string
	order  Order
}

func NewOrderValidator(order Order) *orderValidator {
	return &orderValidator{
		errors: make(map[string]string),
		order:  order,
	}
}

func (v *orderValidator) Validate() map[string]string {
	if v.order.CustomerName == "" {
		v.errors["CustomerName"] = "Customer name is required"
	}
	if len(v.order.Items) < 1 {
		v.errors["Items"] = "At least one order item is required"
	}

	menuIDSet := make(map[int]bool)
	for _, item := range v.order.Items {
		key := "Items[" + strconv.Itoa(item.MenuID) + "]"

		if menuIDSet[item.MenuID] {
			v.errors[key+".MenuID"] = "Duplicate menu item ID detected"
		} else {
			menuIDSet[item.MenuID] = true
		}

		if item.Quantity < 1 {
			v.errors[key+".Quantity"] = "Quantity must be 1 or more"
		}
	}

	if len(v.errors) > 0 {
		return v.errors
	}
	return nil
}
