package repository

import (
	"database/sql"
	"frappuccino/internal/models"
	"log/slog"

	"github.com/lib/pq"
)

type MenuRepository interface {
	InsertMenuItem(tx *sql.Tx, item models.MenuItem) (int, error)
	InsertMenuInventory(tx *sql.Tx, menuID int, inventory []models.MenuItemInventory) error
	RetrieveAll()
	RetrieveByID()
	Update()
	Delete()
	BeginTransaction() (*sql.Tx, error)
}

type menuRepositoryPostgres struct {
	pq     *sql.DB
	logger *slog.Logger
}

func NewMenuRepositoryPostgres(db *sql.DB, logger *slog.Logger) *menuRepositoryPostgres {
	return &menuRepositoryPostgres{
		pq:     db,
		logger: logger,
	}
}

func (m *menuRepositoryPostgres) BeginTransaction() (*sql.Tx, error) {
	tx, err := m.pq.Begin()

	return tx, err
}

func (m *menuRepositoryPostgres) InsertMenuItem(tx *sql.Tx, menuItem models.MenuItem) (int, error) {
	stmt, err := m.pq.Prepare("INSERT INTO menu_items (name, description, price) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		m.logger.Error("Failed to prepare statement", "error", err)
		return -1, err
	}
	defer stmt.Close()

	var menuID int
	err = tx.Stmt(stmt).QueryRow(menuItem.Name, menuItem.Description, menuItem.Price).Scan(&menuID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return -1, models.ErrDuplicateMenuItem
			case "23514":
				return -1, models.ErrNegativePrice
			}
		}
		return -1, err
	}

	return menuID, nil
}

func (m *menuRepositoryPostgres) InsertMenuInventory(tx *sql.Tx, menuID int, inventory []models.MenuItemInventory) error {
	stmt, err := m.pq.Prepare("INSERT INTO menu_item_inventory (menu_id, inventory_id, quantity) VALUES($1, $2, $3)")
	if err != nil {
		m.logger.Error("Failed to prepare statement", "error", err)
		return err
	}
	defer stmt.Close()

	for _, inv := range inventory {
		_, err = tx.Stmt(stmt).Exec(menuID, inv.InventoryID, inv.Quantity)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok {
				switch pgErr.Code {
				case "23503":
					return models.ErrForeignKeyConstraintMenuInventory
				case "23514":
					return models.ErrNegativeQuantity
				}
			}
			return err
		}
	}

	return nil
}

func (m *menuRepositoryPostgres) RetrieveAll() {

}

func (m *menuRepositoryPostgres) RetrieveByID() {

}

func (m *menuRepositoryPostgres) Update() {

}

func (m *menuRepositoryPostgres) Delete() {

}
