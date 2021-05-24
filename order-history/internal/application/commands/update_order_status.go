package commands

import (
	"context"

	"github.com/stackus/ftgogo/order-history/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type UpdateOrderStatus struct {
	OrderID string
	Status  orderapi.OrderState
}

type UpdateOrderStatusHandler struct {
	repo domain.OrderHistoryRepository
}

func NewUpdateOrderStatusHandler(orderHistoryRepo domain.OrderHistoryRepository) UpdateOrderStatusHandler {
	return UpdateOrderStatusHandler{repo: orderHistoryRepo}
}

func (h UpdateOrderStatusHandler) Handle(ctx context.Context, cmd UpdateOrderStatus) error {
	return h.repo.UpdateStatus(ctx, cmd.OrderID, cmd.Status)
}
