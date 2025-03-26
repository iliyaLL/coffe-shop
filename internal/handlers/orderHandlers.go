package handlers

import (
	"encoding/json"
	"fmt"
	"frappuccino/internal/models"
	"frappuccino/internal/utils"
	"net/http"
)

func (app *application) orderCreate(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": "request body does not match json format"})
		return
	}
	defer r.Body.Close()

	m, err := app.OrderSvc.Insert(order)
	if err != nil {
		status, body := mapErrorToResponse(err, m)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, utils.Response{"message": "created"})
}

func (app *application) orderRetrieveAll(w http.ResponseWriter, r *http.Request) {
	orders, err := app.OrderSvc.RetrieveAll()
	if err != nil {
		status, body := mapErrorToResponse(err, nil)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, orders)
}

func (app *application) orderRetrieveByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := app.OrderSvc.RetrieveByID(id)
	if err != nil {
		status, body := mapErrorToResponse(err, nil)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, order)
}

func (app *application) orderUpdateByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, utils.Response{"error": "request body does not match json format"})
		return
	}
	defer r.Body.Close()

	m, err := app.OrderSvc.Update(id, order)
	if err != nil {
		status, body := mapErrorToResponse(err, m)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Response{"message": "OK"})
}

func (app *application) orderDeleteByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := app.OrderSvc.Delete(id)
	if err != nil {
		status, body := mapErrorToResponse(err, nil)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Response{"message": fmt.Sprintf("Deleted %s", id)})
}

func (app *application) orderCloseByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := app.OrderSvc.Close(id)
	if err != nil {
		status, body := mapErrorToResponse(err, nil)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Response{"message": fmt.Sprintf("Closed %s", id)})
}
