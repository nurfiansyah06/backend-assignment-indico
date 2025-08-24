package handler

import "github.com/gin-gonic/gin"

type OrderHandler interface {
	CreateOrder(c *gin.Context)
	GetOrderDetail(c *gin.Context)
}
