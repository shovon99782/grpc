package service

import (
	"context"
	"database/sql"
	"errors"
)

type StockService struct {
	db *sql.DB
}

func NewStockService(db *sql.DB) *StockService {
	return &StockService{db: db}
}

func (s *StockService) ReserveStock(orderID string, items map[string]int) error {
	tx, err := s.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for sku, qty := range items {
		var available int

		// Lock row (FOR UPDATE) → prevents race conditions
		err := tx.QueryRow(`
			SELECT qty_available FROM stocks WHERE sku = ? FOR UPDATE
		`, sku).Scan(&available)

		if err == sql.ErrNoRows {
			return errors.New("SKU not found: " + sku)
		}
		if err != nil {
			return err
		}

		if available < qty {
			return errors.New("not enough stock for SKU: " + sku)
		}

		// Update stock
		_, err = tx.Exec(`
			UPDATE stocks 
			SET qty_available = qty_available - ?, qty_reserved = qty_reserved + ?
			WHERE sku = ?
		`, qty, qty, sku)
		if err != nil {
			return err
		}

		// Insert reservation record
		_, err = tx.Exec(`
			INSERT INTO reservations (order_id, sku, quantity)
			VALUES (?, ?, ?)
		`, orderID, sku, qty)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *StockService) ReleaseStock(orderID string, items map[string]int) error {
	tx, err := s.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for sku, qty := range items {

		// Lock SKU row so race conditions won't happen
		var reserved int
		err := tx.QueryRow(`
			SELECT qty_reserved FROM stocks WHERE sku = ? FOR UPDATE
		`, sku).Scan(&reserved)

		if err == sql.ErrNoRows {
			return errors.New("SKU not found: " + sku)
		}
		if err != nil {
			return err
		}

		if reserved < qty {
			return errors.New("not enough reserved stock to release for SKU: " + sku)
		}

		// Release stock (increase available, decrease reserved)
		_, err = tx.Exec(`
			UPDATE stocks
			SET qty_available = qty_available + ?, qty_reserved = qty_reserved - ?
			WHERE sku = ?
		`, qty, qty, sku)
		if err != nil {
			return err
		}

		// Mark reservation as released
		_, err = tx.Exec(`
			UPDATE reservations
			SET released = TRUE
			WHERE order_id = ? AND sku = ? AND released = FALSE
		`, orderID, sku)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *StockService) GetStock(sku string) (int, int, error) {
	var available, reserved int

	err := s.db.QueryRow(`
		SELECT qty_available, qty_reserved 
		FROM stocks WHERE sku = ?
	`, sku).Scan(&available, &reserved)

	if err == sql.ErrNoRows {
		return 0, 0, errors.New("SKU not found")
	}
	if err != nil {
		return 0, 0, err
	}

	return available, reserved, nil
}
