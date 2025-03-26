package service

import (
	"database/sql"
	"frappuccino/internal/models"
	"frappuccino/internal/repository"
	"log/slog"
)

type ReportService interface {
	GetTotalSales() (models.ReportTotalSales, error)
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
