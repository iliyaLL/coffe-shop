package repository

import (
	"database/sql"
	"frappuccino/internal/models"
	"log/slog"
)

type ReportRepository interface {
	GetTotalSales() (models.ReportTotalSales, error)
	GetPopularMenuItems() ([]models.ReportPopularItem, error)
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

func (m *reportRepositoryPostgres) GetPopularMenuItems() ([]models.ReportPopularItem, error) {
	query := `
		SELECT
			mi.name,
			mi.description,
			mi.price,
			SUM(oi.quantity) AS total_items_sold,
			RANK() OVER (ORDER BY SUM(oi.quantity) DESC) AS rank
		FROM order_item oi
		JOIN menu_items mi ON oi.menu_item_id = mi.id
		JOIN orders o ON oi.order_id = o.id
		WHERE o.order_status = 'closed'
		GROUP BY mi.id
		ORDER BY total_items_sold DESC
		LIMIT 5;
	`
	rows, err := m.pq.Query(query)
	if err != nil {
		m.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var popularItems []models.ReportPopularItem
	for rows.Next() {
		var item models.ReportPopularItem

		err = rows.Scan(&item.Name, &item.Description, &item.Price, &item.TotalItemsSold, &item.Rank)
		if err != nil {
			m.logger.Error(err.Error())
			return nil, err
		}

		popularItems = append(popularItems, item)
	}

	return popularItems, nil
}
