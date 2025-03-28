package service

import (
	"database/sql"
	"log/slog"
	"strconv"

	"frappuccino/internal/models"
	"frappuccino/internal/repository"
)

type OrderService interface {
	Insert(order models.Order) (map[string]string, error)
	RetrieveAll() ([]models.Order, error)
	RetrieveByID(id string) (models.Order, error)
	Update(id string, order models.Order) (map[string]string, error)
	Delete(id string) error
	Close(id string) error
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(db *sql.DB, logger *slog.Logger) *orderService {
	return &orderService{
		repository.NewOrderRepositoryPostgres(db, logger),
	}
}

func (s *orderService) Insert(order models.Order) (map[string]string, error) {
	validator := models.NewOrderValidator(order)
	if errMap := validator.Validate(); errMap != nil {
		return errMap, models.ErrMissingFields
	}

	err := s.orderRepo.Insert(order)
	return nil, err
}

func (s *orderService) RetrieveAll() ([]models.Order, error) {
	orders, err := s.orderRepo.RetrieveAll()

	return orders, err
}

func (s *orderService) RetrieveByID(id string) (models.Order, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return models.Order{}, models.ErrInvalidID
	}

	order, err := s.orderRepo.RetrieveByID(idInt)

	return order, err
}

func (s *orderService) Update(id string, order models.Order) (map[string]string, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, models.ErrInvalidID
	}

	validator := models.NewOrderValidator(order)
	if errMap := validator.Validate(); errMap != nil {
		return errMap, models.ErrMissingFields
	}

	return nil, s.orderRepo.Update(idInt, order)
}

func (s *orderService) Delete(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return models.ErrInvalidID
	}

	return s.orderRepo.Delete(idInt)
}

func (s *orderService) Close(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return models.ErrInvalidID
	}

	return s.orderRepo.Close(idInt)
}
