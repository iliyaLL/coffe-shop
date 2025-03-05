package repository

import (
	"database/sql"
	"frappuccino/internal/models"
	"log"
)

type MenuRepository interface {
	InsertMenu(tx *sql.Tx, item models.MenuItem) (int, error)
	InsertMenuInventory(tx *sql.Tx, menuID int, inventory []models.MenuItemInventory) error
	RetrieveAll()
	RetrieveByID()
	Update()
	Delete()
}

type menuRepositoryPostgres struct {
	pq *sql.DB
}

func NewMenuRepositoryPostgres(db *sql.DB) *menuRepositoryPostgres {
	return &menuRepositoryPostgres{pq: db}
}

func (m *menuRepositoryPostgres) InsertMenu(tx *sql.Tx, item models.MenuItem) (int, error) {
	stmt, err := m.pq.Prepare("INSERT INTO menu_items (name, description, price) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var menuID int
	err = tx.Stmt(stmt).QueryRow(item.Name, item.Description, item.Price).Scan(&menuID)

	return menuID, err
}

func (m *menuRepositoryPostgres) InsertMenuInventory(tx *sql.Tx, menuID int, inventory []models.MenuItemInventory) error {
	stmt, err := m.pq.Prepare("INSERT INTO menu_item_ingredient (menu_id, ingredient_id, quantity) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, inv := range inventory {
		_, err = tx.Stmt(stmt).Exec(menuID, inv.InventoryID, inv.Quantity)
		if err != nil {
			// TODO handle error
		}
	}
}

func (m *menuRepositoryPostgres) RetrieveAll() {

}

func (m *menuRepositoryPostgres) RetrieveByID() {

}

func (m *menuRepositoryPostgres) Update() {

}

func (m *menuRepositoryPostgres) Delete() {

}
