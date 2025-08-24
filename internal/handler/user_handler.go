package handler

import (
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(c *gin.Context)
}
