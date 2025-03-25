package handlers

import (
	"errors"
	"frappuccino/internal/models"
	"frappuccino/internal/service"
	"frappuccino/internal/utils"
	"log/slog"
	"net/http"
)

type application struct {
	logger       *slog.Logger
	InventorySvc service.InventoryService
	MenuSvc      service.MenuService
	OrderSvc     service.OrderService
	// add more services
}

func NewApplication(logger *slog.Logger, inventorySvc service.InventoryService, menuSvc service.MenuService, orderSvc service.OrderService) *application {
	return &application{
		logger:       logger,
		InventorySvc: inventorySvc,
		MenuSvc:      menuSvc,
		OrderSvc:     orderSvc,
		// add more services
	}
}

func mapErrorToResponse(err error, validationMap any) (int, any) {
	switch {
	// General errors
	case errors.Is(err, models.ErrInvalidID):
		return http.StatusBadRequest, utils.Response{"error": err.Error()}
	case errors.Is(err, models.ErrNoRecord):
		return http.StatusNotFound, utils.Response{"error": err.Error()}
	case errors.Is(err, models.ErrMissingFields):
		return http.StatusBadRequest, validationMap

	// Inventory errors
	case errors.Is(err, models.ErrDuplicateInventory),
		errors.Is(err, models.ErrNegativeQuantity),
		errors.Is(err, models.ErrInvalidEnumTypeInventory):
		return http.StatusBadRequest, utils.Response{"error": err.Error()}

	// Menu errors
	case errors.Is(err, models.ErrDuplicateMenuItem),
		errors.Is(err, models.ErrNegativePrice),
		errors.Is(err, models.ErrForeignKeyConstraintMenuInventory):
		return http.StatusBadRequest, utils.Response{"error": err.Error()}

	// Order errors
	case errors.Is(err, models.ErrDuplicateOrder),
		errors.Is(err, models.ErrForeignKeyConstraintOrderMenu):
		return http.StatusBadRequest, utils.Response{"error": err.Error()}

	// Default catch-all
	default:
		return http.StatusInternalServerError, utils.Response{"error": "Internal Server Error"}
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

		"POST /menu":        app.menuCreatePost,
		"GET /menu":         app.menuRetrieveAllGet,
		"GET /menu/{id}":    app.menuRetrieveAllByIDGet,
		"PUT /menu/{id}":    app.menuUpdate,
		"DELETE /menu/{id}": app.menuDelete,

		"POST /orders":            app.orderCreate,
		"GET /orders":             app.orderRetrieveAll,
		"GET /orders/{id}":        app.orderRetrieveByID,
		"PUT /orders/{id}":        app.orderUpdateByID,
		"DELETE /orders/{id}":     app.orderDeleteByID,
		"POST /orders/{id}/close": app.orderCloseByID,
	}
	for endpoint, f := range endpoints {
		router.HandleFunc(endpoint, ChainMiddleware(f, commonMiddleware...))
	}

	return router
}
