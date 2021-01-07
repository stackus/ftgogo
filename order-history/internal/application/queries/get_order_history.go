package queries

import (
	"context"

	"github.com/stackus/ftgogo/order-history/internal/domain"
)

type GetOrderHistory struct {
	OrderID string
}

type GetOrderHistoryHandler struct {
	repo domain.OrderHistoryRepository
}

func NewGetOrderHistoryHandler(orderHistoryRepo domain.OrderHistoryRepository) GetOrderHistoryHandler {
	return GetOrderHistoryHandler{repo: orderHistoryRepo}
}

func (h GetOrderHistoryHandler) Handle(ctx context.Context, query GetOrderHistory) (OrderHistory, error) {
	order, err := h.repo.Find(ctx, query.OrderID)
	if err != nil {
		return OrderHistory{}, err
	}

	return OrderHistory{
		OrderID:        order.OrderID,
		Status:         order.Status.String(),
		RestaurantID:   order.RestaurantID,
		RestaurantName: order.RestaurantName,
		CreatedAt:      order.CreatedAt,
	}, nil
}
