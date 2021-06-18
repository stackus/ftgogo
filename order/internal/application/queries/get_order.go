package queries

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type GetOrder struct {
	OrderID string
}

type GetOrderHandler struct {
	orderRepo adapters.OrderRepository
}

func NewGetOrderHandler(orderRepo adapters.OrderRepository) GetOrderHandler {
	return GetOrderHandler{orderRepo: orderRepo}
}

func (h GetOrderHandler) Handle(ctx context.Context, query GetOrder) (*domain.Order, error) {
	return h.orderRepo.Load(ctx, query.OrderID)
}
