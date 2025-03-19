package service

import (
	"database/sql"
	"frappuccino/internal/models"
	"frappuccino/internal/repository"
	"log/slog"
	"strconv"
)

type MenuService interface {
	InsertMenu(menu models.MenuItem) (map[string]string, error)
	RetrieveAll() ([]models.MenuItem, error)
	RetrieveByID(id string) (models.MenuItem, error)
	Update()
	Delete()
}

type menuService struct {
	menuRepo repository.MenuRepository
}

func NewMenuService(db *sql.DB, logger *slog.Logger) *menuService {
	return &menuService{
		repository.NewMenuRepositoryPostgres(db, logger),
	}
}

func (s *menuService) InsertMenu(menu models.MenuItem) (map[string]string, error) {
	validator := models.NewMenuItemValidator(menu)
	if errMap := validator.Validate(); errMap != nil {
		return errMap, models.ErrMissingFields
	}

	tx, err := s.menuRepo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	menuID, err := s.menuRepo.InsertMenuItem(tx, menu)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.menuRepo.InsertMenuInventory(tx, menuID, menu.Inventory)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return nil, tx.Commit()
}

func (s *menuService) RetrieveAll() ([]models.MenuItem, error) {
	menuItems, err := s.menuRepo.RetrieveAll()

	return menuItems, err
}

func (s *menuService) RetrieveByID(id string) (models.MenuItem, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return models.MenuItem{}, models.ErrInvalidID
	}

	menuItem, err := s.menuRepo.RetrieveByID(idInt)

	return menuItem, err
}
func (s *menuService) Update() {}
func (s *menuService) Delete() {}
