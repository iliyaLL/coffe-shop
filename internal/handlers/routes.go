package handlers

import (
	"frappuccino/internal/service"
	"log/slog"
	"net/http"
)

type application struct {
	logger       *slog.Logger
	InventorySvc service.InventoryService
	MenuSvc      service.MenuService
	// add more services
}

func NewApplication(logger *slog.Logger, inventorySvc service.InventoryService, menuSvc service.MenuService) *application {
	return &application{
		logger:       logger,
		InventorySvc: inventorySvc,
		MenuSvc:      menuSvc,
		// add more services
	}
}

func (app *application) Routes() http.Handler {
	router := http.NewServeMux()
	commonMiddleware := []Middleware{
		app.recoverPanic,
		app.logRequest,
		contentTypeJSON,
	}

	endpoints := map[string]http.HandlerFunc{
		"POST /inventory":        app.inventoryCreatePost,
		"GET /inventory":         app.inventoryRetreiveAllGet,
		"GET /inventory/{id}":    app.inventoryRetrieveByIDGet,
		"PUT /inventory/{id}":    app.inventoryUpdateByIDPut,
		"DELETE /inventory/{id}": app.inventoryDeleteByIDDelete,

		"POST /menu": app.menuCreatePost,
	}
	for endpoint, f := range endpoints {
		router.HandleFunc(endpoint, ChainMiddleware(f, commonMiddleware...))
	}

	return router
}
