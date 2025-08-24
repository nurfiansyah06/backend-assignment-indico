package settlementhandler

import (
	"backend-assignment/constant"
	"backend-assignment/internal/handler"
	"backend-assignment/internal/models"
	"backend-assignment/internal/service"
	"backend-assignment/logs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SettlementHandlerImpl struct {
	SettlementService service.SettlementService
	Validate          *validator.Validate
	jobQueue          chan models.Job
}

func NewSettlementHandler(settlementService service.SettlementService, validate *validator.Validate, jobQueue chan models.Job) handler.SettlementHandler {
	return &SettlementHandlerImpl{
		SettlementService: settlementService,
		Validate:          validate,
		jobQueue:          jobQueue,
	}
}

// CreateJob implements handler.SettlementHandler.
func (s *SettlementHandlerImpl) CreateJob(c *gin.Context) {
	var reqBody models.RequestJob

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: constant.InvalidRequest,
		})
		return
	}

	if err := s.Validate.Struct(&reqBody); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	job, resp, err := s.SettlementService.CreateJob(ctx, reqBody)
	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	s.jobQueue <- *job

	c.JSON(http.StatusAccepted, resp)
}

// GetJobByID implements handler.SettlementHandler.
func (s *SettlementHandlerImpl) GetJobById(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		logs.Error("job ID is required")
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "job ID is required",
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := s.SettlementService.GetJobById(ctx, jobID)
	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
