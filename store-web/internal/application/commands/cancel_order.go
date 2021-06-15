package commands

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/store-web/internal/adapters"
)

type CancelOrder struct {
	OrderID string
}

type CancelOrderHandler struct {
	repo adapters.OrderRepository
}

func NewCancelOrderHandler(repo adapters.OrderRepository) CancelOrderHandler {
	return CancelOrderHandler{repo: repo}
}

func (h CancelOrderHandler) Handle(ctx context.Context, cmd CancelOrder) (orderapi.OrderState, error) {
	_, err := h.repo.Find(ctx, adapters.FindOrder{
		OrderID: cmd.OrderID,
	})
	if err != nil {
		return orderapi.UnknownOrderState, err
	}

	return h.repo.Cancel(ctx, adapters.CancelOrder{
		OrderID: cmd.OrderID,
	})
}
