package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type ConfirmCancelOrder struct {
	OrderID string
}

type ConfirmCancelOrderHandler struct {
	repo adapters.OrderRepository
}

func NewConfirmCancelOrderHandler(repo adapters.OrderRepository) ConfirmCancelOrderHandler {
	return ConfirmCancelOrderHandler{
		repo: repo,
	}
}

func (h ConfirmCancelOrderHandler) Handle(ctx context.Context, cmd ConfirmCancelOrder) error {
	_, err := h.repo.Update(ctx, cmd.OrderID, &domain.ConfirmCancelOrder{})
	if err != nil {
		return err
	}

	return nil
}
