package handlers

import (
	"encoding/json"
	"errors"
	"frappuccino/internal/models"
	"frappuccino/internal/utils"
	"net/http"
)

func (app *application) inventoryCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var inventory models.Inventory
	err := json.NewDecoder(r.Body).Decode(&inventory)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "request body does not match json format"})
		return
	}
	defer r.Body.Close()

	m, err := app.InventorySvc.Insert(&inventory)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateInventory) {
			utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "duplicate inventory"})
		} else if errors.Is(err, models.ErrNegativeQuantity) {
			utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "negative quantity"})
		} else if errors.Is(err, models.ErrMissingFields) {
			utils.SendJSONResponse(w, http.StatusBadRequest, m)
		} else if errors.Is(err, models.ErrInvalidEnumType) {
			utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid unit type", "supported types": "shots, ml, g, units"})
		} else {
			utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, map[string]string{"message": "created"})
}
