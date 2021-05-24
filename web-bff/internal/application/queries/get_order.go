package queries

import (
	"context"

	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type GetOrder struct {
	OrderID string
}

type GetOrderHandler struct {
	repo domain.OrderRepository
}

func NewGetOrderHandler(repo domain.OrderRepository) GetOrderHandler {
	return GetOrderHandler{repo: repo}
}

func (h GetOrderHandler) Handle(ctx context.Context, query GetOrder) (*domain.Order, error) {
	order, err := h.repo.Find(ctx, query.OrderID)
	return order, err
}
