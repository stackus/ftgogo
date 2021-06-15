package queries

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/adapters"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type GetOrder struct {
	OrderID string
}

type GetOrderHandler struct {
	repo adapters.OrderRepository
}

func NewGetOrderHandler(repo adapters.OrderRepository) GetOrderHandler {
	return GetOrderHandler{repo: repo}
}

func (h GetOrderHandler) Handle(ctx context.Context, query GetOrder) (*domain.Order, error) {
	return h.repo.Find(ctx, adapters.FindOrder{
		OrderID: query.OrderID,
	})
}
