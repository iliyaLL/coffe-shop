package repository

import (
	"database/sql"
	"frappuccino/internal/models"
	"log/slog"
)

type ReportRepository interface {
	GetTotalSales() (models.ReportTotalSales, error)
}

type reportRepositoryPostgres struct {
	pq     *sql.DB
	logger *slog.Logger
}

func NewReportRepositoryPostgres(db *sql.DB, logger *slog.Logger) *reportRepositoryPostgres {
	return &reportRepositoryPostgres{
		pq:     db,
		logger: logger,
	}
}

func (m *reportRepositoryPostgres) GetTotalSales() (models.ReportTotalSales, error) {
	query := `
		SELECT
			COUNT(DISTINCT o.id) AS orders_completed,
			COALESCE(SUM(oi.quantity * mi.price), 0) AS total_sales
		FROM orders o
		JOIN order_item oi ON o.id = oi.order_id
		JOIN menu_items mi ON oi.menu_item_id = mi.id
		WHERE o.order_status = 'closed';
	`
	var report models.ReportTotalSales
	err := m.pq.QueryRow(query).Scan(&report.OrdersCompleted, &report.TotalSales)
	if err != nil {
		m.logger.Error(err.Error())
		return models.ReportTotalSales{}, err
	}

	return report, nil
}
