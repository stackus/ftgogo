package commands

import (
	"context"

	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type CancelOrder struct {
	ConsumerID string
	OrderID    string
}

type CancelOrderHandler struct {
	repo domain.OrderRepository
}

func NewCancelOrderHandler(repo domain.OrderRepository) CancelOrderHandler {
	return CancelOrderHandler{repo: repo}
}

func (h CancelOrderHandler) Handle(ctx context.Context, cmd CancelOrder) (orderapi.OrderState, error) {
	order, err := h.repo.Find(ctx, domain.FindOrder{
		OrderID: cmd.OrderID,
	})
	if err != nil {
		return orderapi.UnknownOrderState, err
	}

	if order.ConsumerID != cmd.ConsumerID {
		// being opaque intentionally; Could also send a permission denied error
		return orderapi.UnknownOrderState, errors.Wrap(errors.ErrNotFound, "order not found")
	}

	return h.repo.Cancel(ctx, domain.CancelOrder{
		OrderID: cmd.OrderID,
	})
}
