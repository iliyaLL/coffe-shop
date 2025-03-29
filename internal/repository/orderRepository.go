package repository

import (
	"database/sql"
	"encoding/json"
	"frappuccino/internal/models"
	"log/slog"
	"time"

	"github.com/lib/pq"
)

type OrderRepository interface {
	Insert(order models.Order) error
	RetrieveAll() ([]models.Order, error)
	RetrieveByID(id int) (models.Order, error)
	Update(orderID int, order models.Order) error
	Delete(id int) error
	Close(id int) error
}

type orderRepositoryPostgres struct {
	pq     *sql.DB
	logger *slog.Logger
}

func NewOrderRepositoryPostgres(db *sql.DB, logger *slog.Logger) *orderRepositoryPostgres {
	return &orderRepositoryPostgres{
		pq:     db,
		logger: logger,
	}
}

func (m *orderRepositoryPostgres) Insert(order models.Order) error {
	tx, err := m.pq.Begin()
	if err != nil {
		m.logger.Error("Failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback()

	prefsJSON, err := json.Marshal(order.CustomerPreferences)
	if err != nil {
		m.logger.Error("Failed to marshal customer_preferences", "error", err)
		return err
	}

	var orderID int
	err = tx.QueryRow(`INSERT INTO orders (customer_name, order_status, customer_preferences) 
		VALUES ($1, $2, $3) RETURNING id`,
		order.CustomerName, "open", prefsJSON).
		Scan(&orderID)
	if err != nil {
		m.logger.Error(err.Error())
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return models.ErrDuplicateOrder
			case "22P02":
				return models.ErrInvalidEnumTypeInventory
			}
		}
		return err
	}

	for _, menu := range order.Items {
		_, err = tx.Exec("INSERT INTO order_item (order_id, menu_item_id, quantity) VALUES ($1, $2, $3)",
			orderID, menu.MenuID, menu.Quantity)
		if err != nil {
			m.logger.Error(err.Error())
			if pgErr, ok := err.(*pq.Error); ok {
				switch pgErr.Code {
				case "23503":
					return models.ErrForeignKeyConstraintOrderMenu
				case "23514":
					return models.ErrNegativeQuantity
				}
			}
			return err
		}

		rows, err := tx.Query("SELECT inventory_id, quantity FROM menu_item_inventory WHERE menu_id=$1", menu.MenuID)
		if err != nil {
			m.logger.Error(err.Error())
			return err
		}
		defer rows.Close()

		type invUse struct {
			inventoryID int
			totalNeeded int
		}
		var inventoryList []invUse
		for rows.Next() {
			var inventoryID, perItemQuantity int
			if err := rows.Scan(&inventoryID, &perItemQuantity); err != nil {
				return err
			}
			inventoryList = append(inventoryList, invUse{inventoryID, perItemQuantity * menu.Quantity})
		}

		for _, item := range inventoryList {
			_, err := tx.Exec("UPDATE inventory SET quantity = quantity - $1 WHERE id = $2", item.totalNeeded, item.inventoryID)
			if err != nil {
				m.logger.Error(err.Error())
				if pqErr, ok := err.(*pq.Error); ok {
					switch pqErr.Code {
					case "23514":
						return models.ErrNegativeQuantity
					}
				}
				return err
			}
		}
	}

	return tx.Commit()
}

func (m *orderRepositoryPostgres) RetrieveAll() ([]models.Order, error) {
	rows, err := m.pq.Query(`
		SELECT o.id, o.customer_name, o.order_status, o.created_at, o.customer_preferences,
		       oi.menu_item_id, oi.quantity
		FROM orders o
		LEFT JOIN order_item oi ON o.id = oi.order_id
		ORDER BY o.id
	`)
	if err != nil {
		m.logger.Error("Failed to execute order query", "error", err)
		return nil, err
	}
	defer rows.Close()

	orderMap := make(map[int]*models.Order)

	for rows.Next() {
		var (
			orderID      int
			customerName string
			status       string
			createdAt    time.Time
			prefsBytes   []byte
			menuItemID   sql.NullInt32
			quantity     sql.NullInt32
		)

		err := rows.Scan(&orderID, &customerName, &status, &createdAt, &prefsBytes, &menuItemID, &quantity)
		if err != nil {
			m.logger.Error("Failed to scan order row", "error", err)
			return nil, err
		}

		if _, ok := orderMap[orderID]; !ok {
			var prefs models.Jsonb
			if err := json.Unmarshal(prefsBytes, &prefs); err != nil {
				m.logger.Error("Failed to unmarshal customer_preferences", "error", err)
				return nil, err
			}

			orderMap[orderID] = &models.Order{
				ID:                  orderID,
				CustomerName:        customerName,
				Status:              status,
				CreatedAt:           createdAt,
				CustomerPreferences: prefs,
				Items:               []models.OrderItem{},
			}
		}

		if menuItemID.Valid {
			orderMap[orderID].Items = append(orderMap[orderID].Items, models.OrderItem{
				MenuID:   int(menuItemID.Int32),
				Quantity: int(quantity.Int32),
			})
		}
	}

	var orders []models.Order
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (m *orderRepositoryPostgres) RetrieveByID(id int) (models.Order, error) {
	rows, err := m.pq.Query(`
		SELECT o.id, o.customer_name, o.order_status, o.created_at, o.customer_preferences,
		       oi.menu_item_id, oi.quantity
		FROM "orders" o
		LEFT JOIN order_item oi ON o.id = oi.order_id
		WHERE o.id = $1
	`, id)
	if err != nil {
		m.logger.Error("Failed to execute order query", "error", err)
		return models.Order{}, err
	}
	defer rows.Close()

	var order models.Order
	for rows.Next() {
		var (
			orderID      int
			customerName string
			status       string
			createdAt    time.Time
			prefsBytes   []byte
			menuItemID   sql.NullInt32
			quantity     sql.NullInt32
		)

		err := rows.Scan(&orderID, &customerName, &status, &createdAt, &prefsBytes, &menuItemID, &quantity)
		if err != nil {
			m.logger.Error("Failed to scan order row", "error", err)
			return models.Order{}, err
		}

		if order.ID == 0 {
			var prefs models.Jsonb
			if err := json.Unmarshal(prefsBytes, &prefs); err != nil {
				m.logger.Error("Failed to unmarshal customer_preferences", "error", err)
				return models.Order{}, err
			}

			order = models.Order{
				ID:                  orderID,
				CustomerName:        customerName,
				Status:              status,
				CreatedAt:           createdAt,
				CustomerPreferences: prefs,
				Items:               []models.OrderItem{},
			}
		}

		if menuItemID.Valid {
			order.Items = append(order.Items, models.OrderItem{
				MenuID:   int(menuItemID.Int32),
				Quantity: int(quantity.Int32),
			})
		}
	}

	if order.ID == 0 {
		return models.Order{}, models.ErrNoRecord
	}

	return order, nil
}

func (m *orderRepositoryPostgres) Update(orderID int, order models.Order) error {
	tx, err := m.pq.Begin()
	if err != nil {
		m.logger.Error("Failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback()

	prefsJSON, err := json.Marshal(order.CustomerPreferences)
	if err != nil {
		m.logger.Error("Failed to marshal customer_preferences", "error", err)
		return err
	}

	result, err := tx.Exec(`
		UPDATE orders
		SET customer_name = $1, order_status = $2, customer_preferences = $3
		WHERE id = $4
	`, order.CustomerName, order.Status, prefsJSON, orderID)
	if err != nil {
		m.logger.Error(err.Error())
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return models.ErrDuplicateOrder
			}
		}
		m.logger.Error("Failed to update order", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		m.logger.Error("Failed to check rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	_, err = tx.Exec("DELETE FROM order_item WHERE order_id = $1", orderID)
	if err != nil {
		m.logger.Error("Failed to delete order items", "error", err)
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(
			"INSERT INTO order_item (order_id, menu_item_id, quantity) VALUES ($1, $2, $3)",
			orderID, item.MenuID, item.Quantity,
		)
		if err != nil {
			m.logger.Error(err.Error())
			if pgErr, ok := err.(*pq.Error); ok {
				switch pgErr.Code {
				case "23503":
					return models.ErrForeignKeyConstraintOrderMenu
				case "23514":
					return models.ErrNegativeQuantity
				}
			}
			m.logger.Error("Failed to insert order item", "menu_id", item.MenuID, "error", err)
			return err
		}
	}

	return tx.Commit()
}

func (m *orderRepositoryPostgres) Delete(id int) error {
	result, err := m.pq.Exec("DELETE FROM orders WHERE id=$1", id)
	if err != nil {
		m.logger.Error("Failed to execute query", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	return err
}

func (m *orderRepositoryPostgres) Close(id int) error {
	result, err := m.pq.Exec(`UPDATE orders SET order_status=$1 WHERE id=$2`, "closed", id)
	if err != nil {
		m.logger.Error("Failed to close order", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		m.logger.Error("Failed to check rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	return nil
}
