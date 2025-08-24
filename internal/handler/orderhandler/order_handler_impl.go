package orderhandler

import (
	"backend-assignment/constant"
	"backend-assignment/internal/handler"
	"backend-assignment/internal/models"
	"backend-assignment/internal/service"
	"backend-assignment/logs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderHandlerImpl struct {
	OrderService service.OrderService
	Validate     *validator.Validate
}

func NewOrderHandler(OrderService service.OrderService, validate *validator.Validate) handler.OrderHandler {
	return &OrderHandlerImpl{
		OrderService: OrderService,
		Validate:     validate,
	}
}

// CreateOrder implements handler.OrderHandler.
func (o *OrderHandlerImpl) CreateOrder(c *gin.Context) {
	var reqBody models.RequestOrder

	// Bind JSON request
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: constant.InvalidRequest,
		})
		return
	}

	// Validate
	if err := o.Validate.Struct(&reqBody); err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	resp, err := o.OrderService.CreateOrder(ctx, reqBody)
	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)

}

// GetOrderDetail implements handler.OrderHandler.
func (o *OrderHandlerImpl) GetOrderDetail(c *gin.Context) {
	idParam := c.Param("id")
	orderId, err := strconv.Atoi(idParam)
	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: constant.InvalidId,
		})
		return
	}

	ctx := c.Request.Context()
	detail, err := o.OrderService.GetOrderDetail(ctx, orderId)
	if err != nil {
		logs.Error(err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: constant.ErrGetOrderDetailService,
		})
		return
	}

	c.JSON(http.StatusOK, detail)
}
