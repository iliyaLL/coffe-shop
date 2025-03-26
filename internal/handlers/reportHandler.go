package handlers

import (
	"frappuccino/internal/utils"
	"net/http"
)

func (app *application) getTotalSalesReport(w http.ResponseWriter, r *http.Request) {
	report, err := app.ReportSvc.GetTotalSales()
	if err != nil {
		status, body := utils.MapErrorToResponse(err, nil)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, report)
}

func (app *application) getPopularMenuItems(w http.ResponseWriter, r *http.Request) {
	popularItems, err := app.ReportSvc.GetPopularMenuItems()
	if err != nil {
		status, body := utils.MapErrorToResponse(err, nil)
		utils.SendJSONResponse(w, status, body)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, popularItems)
}
