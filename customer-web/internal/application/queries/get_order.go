package queries

import (
	"context"

	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/customer-web/internal/adapters"
	"github.com/stackus/ftgogo/customer-web/internal/domain"
)

type GetOrder struct {
	OrderID    string
	ConsumerID string
}

type GetOrderHandler struct {
	repo adapters.OrderRepository
}

func NewGetOrderHandler(repo adapters.OrderRepository) GetOrderHandler {
	return GetOrderHandler{repo: repo}
}

func (h GetOrderHandler) Handle(ctx context.Context, query GetOrder) (*domain.Order, error) {
	order, err := h.repo.Find(ctx, adapters.FindOrder{
		OrderID: query.OrderID,
	})
	if err != nil {
		return nil, err
	}

	if order.ConsumerID != query.ConsumerID {
		// being opaque intentionally; Could also send a permission denied error
		return nil, errors.Wrap(errors.ErrNotFound, "order not found")
	}

	return order, nil
}
