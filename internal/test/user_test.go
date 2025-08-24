package test

import (
	"backend-assignment/internal/models"
	"backend-assignment/internal/service/userservice"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTblUser struct {
	mock.Mock
}

func (m *MockTblUser) Register(ctx context.Context, request models.RequestUser) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockTblUser) GetUser(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func TestRegisterUsersSuccess(t *testing.T) {
	mockRepo := new(MockTblUser)
	userService := userservice.NewUserService(mockRepo)
	req := models.RequestUser{Name: "Mas"}

	mockRepo.On("Register", mock.Anything, req).Return(nil)

	resp, err := userService.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.Status)
	assert.Equal(t, "Mas", resp.Data.(models.RequestUser).Name)
	assert.Equal(t, "Success Registered", resp.Message)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUsersFailed(t *testing.T) {
	mockRepo := new(MockTblUser)
	userService := userservice.NewUserService(mockRepo)

	req := models.RequestUser{Name: "Mas"}
	mockRepo.On("Register", mock.Anything, req).Return(errors.New("db error"))

	resp, err := userService.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestGetUserSuccess(t *testing.T) {
	mockRepo := new(MockTblUser)

	expectedUser := &models.User{Id: 1, Name: "Mas"}
	mockRepo.On("GetUser", mock.Anything, 1).Return(expectedUser, nil)

	user, err := mockRepo.GetUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.Equal(t, expectedUser.Name, user.Name)

	mockRepo.AssertExpectations(t)
}

func TestGetUserNotFound(t *testing.T) {
	mockRepo := new(MockTblUser)

	mockRepo.On("GetUser", mock.Anything, 2).Return(nil, errors.New("user not found"))

	user, err := mockRepo.GetUser(context.Background(), 2)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}
