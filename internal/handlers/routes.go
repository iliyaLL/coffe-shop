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
	OrderSvc     service.OrderService
	ReportSvc    service.ReportService
	// add more services
}

func NewApplication(logger *slog.Logger,
	inventorySvc service.InventoryService,
	menuSvc service.MenuService,
	orderSvc service.OrderService,
	reportSvc service.ReportService,
) *application {
	return &application{
		logger:       logger,
		InventorySvc: inventorySvc,
		MenuSvc:      menuSvc,
		OrderSvc:     orderSvc,
		ReportSvc:    reportSvc,
		// add more services
	}
}

func (app *application) Routes() http.Handler {
	router := http.NewServeMux()
	commonMiddleware := []Middleware{
		app.recoverPanic,
		app.logRequest,
	}

	endpoints := map[string]http.HandlerFunc{
		// inventory endpoints
		"POST /inventory":        app.inventoryCreate,
		"GET /inventory":         app.inventoryRetreiveAll,
		"GET /inventory/{id}":    app.inventoryRetrieveByID,
		"PUT /inventory/{id}":    app.inventoryUpdateByID,
		"DELETE /inventory/{id}": app.inventoryDeleteByID,

		// menu endpoints
		"POST /menu":        app.menuCreate,
		"GET /menu":         app.menuRetrieveAll,
		"GET /menu/{id}":    app.menuRetrieveAllByID,
		"PUT /menu/{id}":    app.menuUpdate,
		"DELETE /menu/{id}": app.menuDelete,

		// orders endpoints
		"POST /orders":            app.orderCreate,
		"GET /orders":             app.orderRetrieveAll,
		"GET /orders/{id}":        app.orderRetrieveByID,
		"PUT /orders/{id}":        app.orderUpdateByID,
		"DELETE /orders/{id}":     app.orderDeleteByID,
		"POST /orders/{id}/close": app.orderCloseByID,

		// new endpoints
		"GET /orders/numberOfOrderedItems": app.numberOfOrderedItems, // ?startDate={startDate}&endDate={endDate}

		// aggregations endpoints
		"GET /reports/total-sales":   app.getTotalSalesReport,
		"GET /reports/popular-items": app.getPopularMenuItems,
	}

	for endpoint, f := range endpoints {
		router.HandleFunc(endpoint, ChainMiddleware(f, commonMiddleware...))
	}

	return router
}
