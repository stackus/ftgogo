package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type BeginCancelOrder struct {
	OrderID string
}

type BeginCancelOrderHandler struct {
	repo ports.OrderRepository
}

func NewBeginCancelOrderHandler(orderRepo ports.OrderRepository) BeginCancelOrderHandler {
	return BeginCancelOrderHandler{repo: orderRepo}
}

func (h BeginCancelOrderHandler) Handle(ctx context.Context, cmd BeginCancelOrder) error {
	_, err := h.repo.Update(ctx, cmd.OrderID, &domain.BeginCancelOrder{})

	return err
}
