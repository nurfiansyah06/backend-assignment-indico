package tblorder

import (
	"backend-assignment/internal/models"
	"backend-assignment/internal/repository/database"
	"backend-assignment/logs"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type TblOrderImpl struct {
	DB *sql.DB
}

func New(DB *sql.DB) database.TblOrder {
	return &TblOrderImpl{
		DB: DB,
	}
}

// CreateOrder implements database.TblOrder.
func (repository *TblOrderImpl) CreateOrder(ctx context.Context, request *models.RequestOrder) (*models.Order, error) {
	tx, err := repository.DB.Begin()
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	var stock int
	err = tx.QueryRowContext(ctx, "SELECT stock FROM products WHERE id=$1 FOR UPDATE", request.ProductId).
		Scan(&stock)
	if err != nil {
		return nil, err
	}

	if stock < request.Quantity {
		return nil, fmt.Errorf("insufficient stock: available %d, requested %d", stock, request.Quantity)
	}

	_, err = tx.ExecContext(ctx, "UPDATE products SET stock = stock - $1 WHERE id=$2", request.Quantity, request.ProductId)
	if err != nil {
		return nil, err
	}

	var order models.Order
	// Fixed SQL syntax - removed the stray backtick and concatenated properly
	err = tx.QueryRowContext(ctx,
		`INSERT INTO orders (product_id, user_id, quantity) VALUES ($1, $2, $3) 
         RETURNING id, product_id, user_id, quantity, date_created`,
		request.ProductId, request.UserId, request.Quantity).
		Scan(&order.Id, &order.ProductId, &order.UserId, &order.Quantity, &order.DateCreated)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrderDetail implements database.TblOrder.
func (repository *TblOrderImpl) GetOrderDetail(ctx context.Context, orderId int) (*models.DetailOrder, error) {
	query := `
		SELECT 
			u.name as user_name,
			p.product_name,
			COALESCE(o.quantity, 0) as quantity
			FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN products p ON o.product_id = p.id
		WHERE o.id = $1`

	var detail models.DetailOrder
	err := repository.DB.QueryRowContext(ctx, query, orderId).Scan(
		&detail.Name,
		&detail.ProductName,
		&detail.Quantity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order detail: %w", err)
	}

	return &detail, nil
}
