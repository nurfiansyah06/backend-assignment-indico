package handler

import "github.com/gin-gonic/gin"

type SettlementHandler interface {
	CreateJob(c *gin.Context)
	GetJobById(c *gin.Context)
}
