package userservice

import (
	"backend-assignment/internal/models"
	"backend-assignment/internal/repository/database"
	"backend-assignment/internal/service"
	"backend-assignment/logs"
	"context"
	"net/http"
)

type UserServiceImpl struct {
	tblUser database.TblUser
}

func NewUserService(tblUser database.TblUser) service.UserService {
	return &UserServiceImpl{
		tblUser: tblUser,
	}
}

// Register implements service.UserUsecase.
func (u *UserServiceImpl) Register(ctx context.Context, request models.RequestUser) (*models.Response, error) {
	user := models.RequestUser{
		Name: request.Name,
	}

	err := u.tblUser.Register(ctx, user)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	response := models.Response{
		Status:  http.StatusCreated,
		Data:    user,
		Message: "Success Registered",
	}

	return &response, nil
}
