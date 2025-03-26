package models

type ReportTotalSales struct {
	OrdersCompleted int     `json:"orders_completed"` // Number of completed orders
	TotalSales      float64 `json:"total_sales"`
}
