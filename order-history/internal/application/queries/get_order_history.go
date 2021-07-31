package queries

import (
	"context"

	"github.com/stackus/ftgogo/order-history/internal/adapters"
	"github.com/stackus/ftgogo/order-history/internal/domain"
)

type GetOrderHistory struct {
	OrderID string
}

type GetOrderHistoryHandler struct {
	repo adapters.OrderHistoryRepository
}

func NewGetOrderHistoryHandler(orderHistoryRepo adapters.OrderHistoryRepository) GetOrderHistoryHandler {
	return GetOrderHistoryHandler{repo: orderHistoryRepo}
}

func (h GetOrderHistoryHandler) Handle(ctx context.Context, query GetOrderHistory) (*domain.OrderHistory, error) {
	return h.repo.Find(ctx, query.OrderID)
}
