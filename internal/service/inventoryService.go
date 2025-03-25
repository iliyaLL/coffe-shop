package service

import (
	"database/sql"
	"frappuccino/internal/models"
	"frappuccino/internal/repository"
	"log/slog"
	"strconv"
)

type InventoryService interface {
	Insert(inventory models.Inventory) (map[string]string, error)
	RetrieveByID(id string) (models.Inventory, error)
	RetrieveAll() ([]models.Inventory, error)
	Update(inventory models.Inventory, id string) (map[string]string, error)
	Delete(id string) error
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
}

func NewInventoryService(db *sql.DB, logger *slog.Logger) *inventoryService {
	return &inventoryService{
		inventoryRepo: repository.NewInventoryRepositoryWithPostgres(db, logger),
	}
}

func (s *inventoryService) Insert(inventory models.Inventory) (map[string]string, error) {
	validator := models.NewInventoryValidator(inventory)
	m := validator.Validate()
	if m != nil {
		return m, models.ErrMissingFields
	}

	err := s.inventoryRepo.Insert(inventory.Name, inventory.Unit, inventory.Quantity, inventory.Categories)

	return nil, err
}

func (s *inventoryService) RetrieveByID(id string) (models.Inventory, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return models.Inventory{}, models.ErrInvalidID
	}

	inventory, err := s.inventoryRepo.RetrieveByID(idInt)

	return inventory, err
}

func (s *inventoryService) RetrieveAll() ([]models.Inventory, error) {
	inventory, err := s.inventoryRepo.RetrieveAll()

	return inventory, err
}

func (s *inventoryService) Update(inventory models.Inventory, id string) (map[string]string, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, models.ErrInvalidID
	}

	validator := models.NewInventoryValidator(inventory)
	m := validator.Validate()
	if len(m) > 0 {
		return m, models.ErrMissingFields
	}

	err = s.inventoryRepo.Update(idInt, inventory.Name, inventory.Unit, inventory.Quantity, inventory.Categories)
	return nil, err
}

func (s *inventoryService) Delete(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return models.ErrInvalidID
	}

	err = s.inventoryRepo.Delete(idInt)
	return err
}
