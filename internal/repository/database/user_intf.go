package database

import (
	"backend-assignment/internal/models"
	"context"
)

type User struct {
	Service TblUser
}

type TblUser interface {
	Register(ctx context.Context, request models.RequestUser) error
	GetUser(ctx context.Context, userId int) (*models.User, error)
}
