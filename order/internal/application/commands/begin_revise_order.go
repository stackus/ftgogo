package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type BeginReviseOrder struct {
	OrderID           string
	RevisedQuantities map[string]int
}

type BeginReviseOrderHandler struct {
	repo domain.OrderRepository
}

func NewBeginReviseOrderHandler(repo domain.OrderRepository) BeginReviseOrderHandler {
	return BeginReviseOrderHandler{
		repo: repo,
	}
}

func (h BeginReviseOrderHandler) Handle(ctx context.Context, cmd BeginReviseOrder) (int, error) {
	root, err := h.repo.Update(ctx, cmd.OrderID, &domain.BeginReviseOrder{
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return 0, err
	}

	order := root.Aggregate().(*domain.Order)

	return order.RevisedOrderTotal(order.OrderTotal(), cmd.RevisedQuantities), nil
}
