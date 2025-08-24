package service

import (
	"backend-assignment/internal/models"
	"context"
)

type OrderService interface {
	CreateOrder(ctx context.Context, request models.RequestOrder) (*models.Response, error)
	GetOrderDetail(ctx context.Context, orderId int) (*models.Response, error)
}
