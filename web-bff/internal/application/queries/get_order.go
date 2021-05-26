package queries

import (
	"context"

	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type GetOrder struct {
	OrderID    string
	ConsumerID string
}

type GetOrderHandler struct {
	repo domain.OrderRepository
}

func NewGetOrderHandler(repo domain.OrderRepository) GetOrderHandler {
	return GetOrderHandler{repo: repo}
}

func (h GetOrderHandler) Handle(ctx context.Context, query GetOrder) (*domain.Order, error) {
	order, err := h.repo.Find(ctx, domain.FindOrder{
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
