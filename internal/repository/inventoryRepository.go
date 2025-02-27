package repository

import (
	"database/sql"
	"errors"
	"frappuccino/internal/models"
	"log"

	"github.com/lib/pq"
)

type InventoryRepository interface {
	Insert(name, unit string, quantity int, categories *[]string) error
	RetrieveByID(id uint32) (models.Inventory, error)
	RetrieveAll() (*[]models.Inventory, error)
}

type InventoryRepositoryPostgres struct {
	pq *sql.DB
}

func NewInventoryRepositoryWithPostgres(db *sql.DB) *InventoryRepositoryPostgres {
	return &InventoryRepositoryPostgres{pq: db}
}

func (model *InventoryRepositoryPostgres) Insert(name, unit string, quantity int, categories *[]string) error {
	stmt, err := model.pq.Prepare("INSERT INTO inventory (name, quantity, unit, categories) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, quantity, unit, categories)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return models.ErrDuplicateInventory
			case "23514":
				return models.ErrNegativeQuantity
			}
		}
		return err
	}

	return nil
}

func (model *InventoryRepositoryPostgres) RetrieveByID(id uint32) (models.Inventory, error) {
	stmt, err := model.pq.Prepare("SELECT * FROM inventory WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	var inventory models.Inventory
	err = stmt.QueryRow(id).Scan(
		&inventory.ID,
		&inventory.Name,
		&inventory.Quantity,
		&inventory.Unit,
		&inventory.Categories,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Inventory{}, models.ErrNoRecord
		}
		return models.Inventory{}, err
	}

	return inventory, nil
}

func (model *InventoryRepositoryPostgres) RetrieveAll() (*[]models.Inventory, error) {
	stmt, err := model.pq.Prepare("SELECT * FROM inventory")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var InventoryAll []models.Inventory
	for rows.Next() {
		var inventory models.Inventory

		err = rows.Scan(
			&inventory.ID,
			&inventory.Name,
			&inventory.Quantity,
			&inventory.Unit,
			&inventory.Categories,
		)
		if err != nil {
			return nil, err
		}

		InventoryAll = append(InventoryAll, inventory)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &InventoryAll, err
}
