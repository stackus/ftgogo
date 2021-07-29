package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type ConfirmReviseOrder struct {
	OrderID           string
	RevisedQuantities map[string]int
}

type ConfirmReviseOrderHandler struct {
	repo ports.OrderRepository
}

func NewConfirmReviseOrderHandler(repo ports.OrderRepository) ConfirmReviseOrderHandler {
	return ConfirmReviseOrderHandler{
		repo: repo,
	}
}

func (h ConfirmReviseOrderHandler) Handle(ctx context.Context, cmd ConfirmReviseOrder) error {
	_, err := h.repo.Update(ctx, cmd.OrderID, &domain.ConfirmReviseOrder{
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return err
	}

	return nil
}
