package repository

import (
	"database/sql"
	"errors"
	"frappuccino/internal/models"
	"log"

	"github.com/lib/pq"
)

type InventoryRepository interface {
	Insert(name, unit string, quantity int, categories []string) error
	RetrieveByID(id int) (models.Inventory, error)
	RetrieveAll() (*[]models.Inventory, error)
	Update(id int, name, unit string, quantity int, categories []string) error
	Delete(id int) error
}

type inventoryRepositoryPostgres struct {
	pq *sql.DB
}

func NewInventoryRepositoryWithPostgres(db *sql.DB) *inventoryRepositoryPostgres {
	return &inventoryRepositoryPostgres{pq: db}
}

func (model *inventoryRepositoryPostgres) Insert(name, unit string, quantity int, categories []string) error {
	stmt, err := model.pq.Prepare("INSERT INTO inventory (name, quantity, unit, categories) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, quantity, unit, pq.Array(categories))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return models.ErrDuplicateInventory
			case "23514":
				return models.ErrNegativeQuantity
			case "22P02":
				return models.ErrInvalidEnumType
			}
		}
		return err
	}

	return nil
}

func (model *inventoryRepositoryPostgres) RetrieveByID(id int) (models.Inventory, error) {
	stmt, err := model.pq.Prepare("SELECT * FROM inventory WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var inventory models.Inventory
	err = stmt.QueryRow(id).Scan(
		&inventory.ID,
		&inventory.Name,
		&inventory.Quantity,
		&inventory.Unit,
		pq.Array(&inventory.Categories),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Inventory{}, models.ErrNoRecord
		}
		return models.Inventory{}, err
	}

	return inventory, nil
}

func (model *inventoryRepositoryPostgres) RetrieveAll() (*[]models.Inventory, error) {
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
			pq.Array(&inventory.Categories),
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

func (model *inventoryRepositoryPostgres) Update(id int, name, unit string, quantity int, categories []string) error {
	stmt, err := model.pq.Prepare("UPDATE inventory SET name=$1, unit=$2, quantity=$3, categories=$4 WHERE id=$5")
	if err != nil {
		log.Fatal()
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, unit, quantity, pq.Array(categories), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecord
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return models.ErrDuplicateInventory
			case "23514":
				return models.ErrNegativeQuantity
			case "22P02":
				return models.ErrInvalidEnumType
			}
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	return err
}

func (model *inventoryRepositoryPostgres) Delete(id int) error {
	stmt, err := model.pq.Prepare("DELETE FROM inventory WHERE id=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	return err
}
