package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type UndoReviseOrder struct {
	OrderID string
}

type UndoReviseOrderHandler struct {
	repo adapters.OrderRepository
}

func NewUndoReviseOrderHandler(orderRepo adapters.OrderRepository) UndoReviseOrderHandler {
	return UndoReviseOrderHandler{repo: orderRepo}
}

func (h UndoReviseOrderHandler) Handle(ctx context.Context, cmd UndoReviseOrder) error {
	_, err := h.repo.Update(ctx, cmd.OrderID, &domain.UndoReviseOrder{})

	return err
}
