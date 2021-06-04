package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type UndoCancelOrder struct {
	OrderID string
}

type UndoCancelOrderHandler struct {
	repo adapters.OrderRepository
}

func NewUndoCancelOrderHandler(orderRepo adapters.OrderRepository) UndoCancelOrderHandler {
	return UndoCancelOrderHandler{repo: orderRepo}
}

func (h UndoCancelOrderHandler) Handle(ctx context.Context, cmd UndoCancelOrder) error {
	_, err := h.repo.Update(ctx, cmd.OrderID, &domain.UndoCancelOrder{})

	return err
}
