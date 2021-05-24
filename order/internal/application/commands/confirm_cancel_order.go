package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type ConfirmCancelOrder struct {
	OrderID string
}

type ConfirmCancelOrderHandler struct {
	repo domain.OrderRepository
}

func NewConfirmCancelOrderHandler(repo domain.OrderRepository) ConfirmCancelOrderHandler {
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
