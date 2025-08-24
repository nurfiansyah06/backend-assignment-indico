package userhandler

import (
	"backend-assignment/internal/handler"
	"backend-assignment/internal/models"
	"backend-assignment/internal/service"
	"backend-assignment/logs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandlerImpl struct {
	UserService service.UserService
	Validate    *validator.Validate
}

func NewUserHandler(userService service.UserService, validate *validator.Validate) handler.UserHandler {
	return &UserHandlerImpl{
		UserService: userService,
		Validate:    validate,
	}
}

// Register implements handler.UserHandler.
func (u *UserHandlerImpl) Register(c *gin.Context) {
	var reqBody models.RequestUser

	// Bind JSON request
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate
	if err := u.Validate.Struct(&reqBody); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	resp, err := u.UserService.Register(c.Request.Context(), reqBody)
	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, resp)
}
