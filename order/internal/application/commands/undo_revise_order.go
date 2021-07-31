package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type UndoReviseOrder struct {
	OrderID string
}

type UndoReviseOrderHandler struct {
	repo ports.OrderRepository
}

func NewUndoReviseOrderHandler(orderRepo ports.OrderRepository) UndoReviseOrderHandler {
	return UndoReviseOrderHandler{repo: orderRepo}
}

func (h UndoReviseOrderHandler) Handle(ctx context.Context, cmd UndoReviseOrder) error {
	_, err := h.repo.Update(ctx, cmd.OrderID, &domain.UndoReviseOrder{})

	return err
}
