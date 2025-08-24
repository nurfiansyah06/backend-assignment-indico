package orderservice

import (
	"backend-assignment/constant"
	"backend-assignment/internal/models"
	"backend-assignment/internal/repository/database"
	"backend-assignment/internal/service"
	"backend-assignment/logs"
	"context"
	"fmt"
	"net/http"
)

type OrderServiceImpl struct {
	tblOrder database.TblOrder
	tblUser  database.TblUser
}

func NewOrderService(tblOrder database.TblOrder, tblUser database.TblUser) service.OrderService {
	return &OrderServiceImpl{
		tblOrder: tblOrder,
		tblUser:  tblUser,
	}
}

// CreateOrder implements service.OrderService.
func (o *OrderServiceImpl) CreateOrder(ctx context.Context, request models.RequestOrder) (*models.Response, error) {
	var response models.Response
	resp, err := o.tblOrder.CreateOrder(ctx, &request)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	product, err := o.tblOrder.GetOrderDetail(ctx, resp.Id)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	user, err := o.tblUser.GetUser(ctx, resp.UserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	response.Status = http.StatusCreated
	response.Data = map[string]interface{}{
		"product_name": product.ProductName,
		"full_name":    user.Name,
		"quantity":     resp.Quantity,
	}
	response.Message = constant.SuccessResponse

	logs.Info(fmt.Sprintf("response: %+v", response.Data))

	return &response, nil

}

// GetOrderDetail implements service.OrderService.
func (o *OrderServiceImpl) GetOrderDetail(ctx context.Context, orderId int) (*models.Response, error) {
	var response models.Response

	detail, err := o.tblOrder.GetOrderDetail(ctx, orderId)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	response.Status = http.StatusOK
	response.Data = detail
	response.Message = constant.SuccessResponse

	logs.Info(fmt.Sprintf("response: %+v", response.Data))

	return &response, nil
}
