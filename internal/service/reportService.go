package service

import (
	"database/sql"
	"log/slog"

	"frappuccino/internal/models"
	"frappuccino/internal/repository"
	"frappuccino/internal/utils"
)

type ReportService interface {
	GetTotalSales() (models.ReportTotalSales, error)
	GetPopularMenuItems() ([]models.ReportPopularItem, error)
	TextSearch(query string, filter string, minPriceArg string, maxPriceArg string) (models.ReportSearch, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(db *sql.DB, logger *slog.Logger) *reportService {
	return &reportService{repository.NewReportRepositoryPostgres(db, logger)}
}

func (s *reportService) GetTotalSales() (models.ReportTotalSales, error) {
	report, err := s.reportRepo.GetTotalSales()

	return report, err
}

func (s *reportService) GetPopularMenuItems() ([]models.ReportPopularItem, error) {
	popularItems, err := s.reportRepo.GetPopularMenuItems()

	return popularItems, err
}

func (s *reportService) TextSearch(query string, filter string, minPriceArg string, maxPriceArg string) (models.ReportSearch, error) {
	if filter != "all" && filter != "menu" && filter != "orders" {
		return models.ReportSearch{}, models.ErrInvalidFilterOption
	}
	minPrice, maxPrice, err := utils.ValidatePrices(minPriceArg, maxPriceArg)
	if err != nil {
		return models.ReportSearch{}, err
	}

	var results models.ReportSearch
	results.TotalMatches = 0
	if filter == "menu" || filter == "all" {
		results.MenuResults, err = s.reportRepo.TextSearchMenu(query, minPrice, maxPrice)
		if err != nil {
			return models.ReportSearch{}, err
		}
		results.TotalMatches += len(results.MenuResults)
	}
	if filter == "orders" || filter == "all" {
		results.OrdersResults, err = s.reportRepo.TextSearchOrders(query, minPrice, maxPrice)
		if err != nil {
			return models.ReportSearch{}, err
		}
		results.TotalMatches += len(results.OrdersResults)
	}
	return results, nil
}
