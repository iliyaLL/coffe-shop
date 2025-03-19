package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"frappuccino/internal/models"
	"frappuccino/internal/utils"
	"net/http"
)

func (app *application) inventoryCreatePost(w http.ResponseWriter, r *http.Request) {
	var inventory models.Inventory
	err := json.NewDecoder(r.Body).Decode(&inventory)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": "request body does not match json format"})
		return
	}
	defer r.Body.Close()

	m, err := app.InventorySvc.Insert(inventory)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateInventory) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else if errors.Is(err, models.ErrNegativeQuantity) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else if errors.Is(err, models.ErrMissingFields) {
			utils.SendJSONResponse(w, http.StatusBadRequest, m)
		} else if errors.Is(err, models.ErrInvalidEnumTypeInventory) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else {
			utils.SendJSONResponse(w, http.StatusInternalServerError, utils.Response{"error": "internal server error"})
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, utils.Response{"message": "created"})
}

func (app *application) inventoryRetreiveAllGet(w http.ResponseWriter, r *http.Request) {
	inventory, err := app.InventorySvc.RetrieveAll()
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, utils.Response{"error": "Internal Server Error"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, inventory)
}

func (app *application) inventoryRetrieveByIDGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	inventory, err := app.InventorySvc.RetrieveByID(id)
	if err != nil {
		if errors.Is(err, models.ErrInvalidID) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else if errors.Is(err, models.ErrNoRecord) {
			utils.SendJSONResponse(w, http.StatusNotFound, utils.Response{"error": err.Error()})
		} else {
			utils.SendJSONResponse(w, http.StatusInternalServerError, utils.Response{"error": "Internal Server error"})
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, inventory)
}

func (app *application) inventoryUpdateByIDPut(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, models.ErrMissingFields) {
			utils.SendJSONResponse(w, http.StatusBadRequest, m)
		} else if errors.Is(err, models.ErrInvalidID) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else if errors.Is(err, models.ErrDuplicateInventory) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else if errors.Is(err, models.ErrNegativeQuantity) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else if errors.Is(err, models.ErrInvalidEnumTypeInventory) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else if errors.Is(err, models.ErrNoRecord) {
			utils.SendJSONResponse(w, http.StatusNotFound, utils.Response{"error": err.Error()})
		} else {
			utils.SendJSONResponse(w, http.StatusInternalServerError, utils.Response{"error": "internal server error"})
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Response{"message": "OK"})
}

func (app *application) inventoryDeleteByIDDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := app.InventorySvc.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrInvalidID) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else if errors.Is(err, models.ErrNoRecord) {
			utils.SendJSONResponse(w, http.StatusNotFound, utils.Response{"error": err.Error()})
		} else {
			utils.SendJSONResponse(w, http.StatusInternalServerError, utils.Response{"error": "Internal Server Error"})
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Response{"message": fmt.Sprintf("Deleted %s", id)})
}
