package service

import (
	"database/sql"
	"frappuccino/internal/models"
	"frappuccino/internal/repository"
	"log/slog"
)

type MenuService interface {
	InsertMenu(menu models.MenuItem) (map[string]string, error)
	RetrieveAll()
	RetrieveByID()
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

func (s *menuService) RetrieveAll()  {}
func (s *menuService) RetrieveByID() {}
func (s *menuService) Update()       {}
func (s *menuService) Delete()       {}
