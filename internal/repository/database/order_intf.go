package database

import (
	"backend-assignment/internal/models"
	"context"
)

type Order struct {
	Service TblOrder
}

type TblOrder interface {
	CreateOrder(ctx context.Context, request *models.RequestOrder) (*models.Order, error)
	GetOrderDetail(ctx context.Context, orderId int) (*models.DetailOrder, error)
}
