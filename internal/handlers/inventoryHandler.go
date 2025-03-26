package handlers

import (
	"encoding/json"
	"fmt"
	"frappuccino/internal/models"
	"frappuccino/internal/utils"
	"net/http"
)

func (app *application) inventoryCreate(w http.ResponseWriter, r *http.Request) {
	var inventory models.Inventory
	err := json.NewDecoder(r.Body).Decode(&inventory)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": "request body does not match json format"})
		return
	}
	defer r.Body.Close()

	m, err := app.InventorySvc.Insert(inventory)
	if err != nil {
		status, body := utils.MapErrorToResponse(err, m)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, utils.Response{"message": "created"})
}

func (app *application) inventoryRetreiveAll(w http.ResponseWriter, r *http.Request) {
	inventory, err := app.InventorySvc.RetrieveAll()
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, utils.Response{"error": "Internal Server Error"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, inventory)
}

func (app *application) inventoryRetrieveByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	inventory, err := app.InventorySvc.RetrieveByID(id)
	if err != nil {
		status, body := utils.MapErrorToResponse(err, nil)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, inventory)
}

func (app *application) inventoryUpdateByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var inventory models.Inventory
	err := json.NewDecoder(r.Body).Decode(&inventory)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": "request body does not match json format"})
		return
	}
	defer r.Body.Close()

	m, err := app.InventorySvc.Update(inventory, id)
	if err != nil {
		status, body := utils.MapErrorToResponse(err, m)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Response{"message": fmt.Sprintf("Updated inventory %s", id)})
}

func (app *application) inventoryDeleteByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := app.InventorySvc.Delete(id)
	if err != nil {
		status, body := utils.MapErrorToResponse(err, nil)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Response{"message": fmt.Sprintf("Deleted %s", id)})
}
