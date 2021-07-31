package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type BeginReviseOrder struct {
	OrderID           string
	RevisedQuantities map[string]int
}

type BeginReviseOrderHandler struct {
	repo ports.OrderRepository
}

func NewBeginReviseOrderHandler(repo ports.OrderRepository) BeginReviseOrderHandler {
	return BeginReviseOrderHandler{
		repo: repo,
	}
}

func (h BeginReviseOrderHandler) Handle(ctx context.Context, cmd BeginReviseOrder) (int, error) {
	order, err := h.repo.Update(ctx, cmd.OrderID, &domain.BeginReviseOrder{
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return 0, err
	}

	return order.RevisedOrderTotal(order.OrderTotal(), cmd.RevisedQuantities), nil
}
