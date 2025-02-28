package service

import (
	"database/sql"
	"frappuccino/internal/models"
	"frappuccino/internal/repository"
)

type InventoryService interface {
	Insert(inventory *models.Inventory) (map[string]string, error)
	RetrieveByID(id int32) (models.Inventory, error)
	RetrieveAll() (*[]models.Inventory, error)
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
}

func NewInventoryService(db *sql.DB) *inventoryService {
	svc := &inventoryService{}
	svc.inventoryRepo = repository.NewInventoryRepositoryWithPostgres(db)
	return svc
}

func (s *inventoryService) Insert(inventory *models.Inventory) (map[string]string, error) {
	validator := models.NewInventoryValidator(inventory)
	m := validator.Validate()
	if len(m) > 0 {
		return m, models.ErrMissingFields
	}

	err := s.inventoryRepo.Insert(inventory.Name, inventory.Unit, inventory.Quantity, inventory.Categories)

	return nil, err
}

func (s *inventoryService) RetrieveByID(id int32) (models.Inventory, error) {
	return models.Inventory{}, nil
}

func (s *inventoryService) RetrieveAll() (*[]models.Inventory, error) {
	inventory, err := s.inventoryRepo.RetrieveAll()

	return inventory, err
}
