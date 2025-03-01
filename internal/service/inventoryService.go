package service

import (
	"database/sql"
	"frappuccino/internal/models"
	"frappuccino/internal/repository"
	"strconv"
)

type InventoryService interface {
	Insert(inventory *models.Inventory) (map[string]string, error)
	RetrieveByID(id string) (models.Inventory, error)
	RetrieveAll() (*[]models.Inventory, error)
	Update(inventory *models.Inventory, id string) (map[string]string, error)
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

func (s *inventoryService) RetrieveByID(id string) (models.Inventory, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return models.Inventory{}, models.ErrInvalidID
	}

	inventory, err := s.inventoryRepo.RetrieveByID(idInt)

	return inventory, err
}

func (s *inventoryService) RetrieveAll() (*[]models.Inventory, error) {
	inventory, err := s.inventoryRepo.RetrieveAll()

	return inventory, err
}

func (s *inventoryService) Update(inventory *models.Inventory, id string) (map[string]string, error) {
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
