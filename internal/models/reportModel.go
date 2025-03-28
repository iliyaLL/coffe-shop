package models

type ReportTotalSales struct {
	OrdersCompleted int     `json:"orders_completed"` // Number of completed orders
	TotalSales      float64 `json:"total_sales"`
}

type ReportPopularItem struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Rank           int     `json:"rank"`
	TotalItemsSold int     `json:"total_items_sold"`
}
