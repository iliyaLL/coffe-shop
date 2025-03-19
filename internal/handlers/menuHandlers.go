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

func (app *application) menuRetrieveAllGet(w http.ResponseWriter, r *http.Request) {
	menuItems, err := app.MenuSvc.RetrieveAll()
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, utils.Response{"error": "Internal Server Error"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, menuItems)
}

func (app *application) menuRetrieveAllByIDGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	menuItem, err := app.MenuSvc.RetrieveByID(id)
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

	utils.SendJSONResponse(w, http.StatusOK, menuItem)
}
