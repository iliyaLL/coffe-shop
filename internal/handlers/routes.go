package handlers

import (
	"frappuccino/internal/service"
	"net/http"
)

type application struct {
	InventorySvc service.InventoryService
	// add more services
}

func NewApplication(inventorySvc service.InventoryService) *application {
	return &application{
		InventorySvc: inventorySvc,
		// add more services
	}
}

func (app *application) Routes() http.Handler {
	router := http.NewServeMux()
	commonMiddleware := []Middleware{
		recoverPanic,
		logRequest,
		contentTypeJSON,
	}

	endpoints := map[string]http.HandlerFunc{
		"POST /inventory":        app.inventoryCreatePost,
		"GET /inventory":         app.inventoryRetreiveAllGet,
		"GET /inventory/{id}":    app.inventoryRetrieveByIDGet,
		"PUT /inventory/{id}":    app.inventoryUpdateByIDPut,
		"DELETE /inventory/{id}": app.inventoryDeleteByIDDelete,
	}
	for endpoint, f := range endpoints {
		router.HandleFunc(endpoint, ChainMiddleware(f, commonMiddleware...))
	}

	return router
}
