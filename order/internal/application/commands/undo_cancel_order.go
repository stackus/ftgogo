package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type UndoCancelOrder struct {
	OrderID string
}

type UndoCancelOrderHandler struct {
	repo ports.OrderRepository
}

func NewUndoCancelOrderHandler(orderRepo ports.OrderRepository) UndoCancelOrderHandler {
	return UndoCancelOrderHandler{repo: orderRepo}
}

func (h UndoCancelOrderHandler) Handle(ctx context.Context, cmd UndoCancelOrder) error {
	_, err := h.repo.Update(ctx, cmd.OrderID, &domain.UndoCancelOrder{})

	return err
}
