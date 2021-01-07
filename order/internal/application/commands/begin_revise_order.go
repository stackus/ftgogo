package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
	"serviceapis/commonapi"
)

type BeginReviseOrder struct {
	OrderID           string
	RevisedQuantities commonapi.MenuItemQuantities
}

type BeginReviseOrderHandler struct {
	repo      domain.OrderRepository
	publisher domain.OrderPublisher
}

func NewBeginReviseOrderHandler(orderRepo domain.OrderRepository, orderPublisher domain.OrderPublisher) BeginReviseOrderHandler {
	return BeginReviseOrderHandler{
		repo:      orderRepo,
		publisher: orderPublisher,
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

	return order.RevisedOrderTotal(order.OrderTotal(), cmd.RevisedQuantities), h.publisher.PublishEntityEvents(ctx, root)
}
