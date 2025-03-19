package handlers

import (
	"encoding/json"
	"errors"
	"frappuccino/internal/models"
	"frappuccino/internal/utils"
	"net/http"
)

func (app *application) menuCreatePost(w http.ResponseWriter, r *http.Request) {
	var menuItem models.MenuItem
	err := json.NewDecoder(r.Body).Decode(&menuItem)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": "request body does not match json format"})
		return
	}
	defer r.Body.Close()

	m, err := app.MenuSvc.InsertMenu(menuItem)
	if err != nil {
		if errors.Is(err, models.ErrMissingFields) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": m})
		} else if errors.Is(err, models.ErrDuplicateMenuItem) {
			utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": err.Error()})
		} else if errors.Is(err, models.ErrForeignKeyConstraintMenuInventory) {
			utils.SendJSONResponse(w, http.StatusNotFound, utils.Response{"error": err.Error()})
		} else {
			utils.SendJSONResponse(w, http.StatusInternalServerError, utils.Response{"error": "internal server error"})
		}

		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, utils.Response{"message": "created"})
}
