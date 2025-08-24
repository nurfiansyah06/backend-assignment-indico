package service

import (
	"backend-assignment/internal/models"
	"context"
)

type UserService interface {
	Register(ctx context.Context, request models.RequestUser) (*models.Response, error)
}
