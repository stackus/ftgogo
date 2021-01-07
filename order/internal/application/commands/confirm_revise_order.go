package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
	"serviceapis/commonapi"
)

type ConfirmReviseOrder struct {
	OrderID           string
	RevisedQuantities commonapi.MenuItemQuantities
}

type ConfirmReviseOrderHandler struct {
	repo      domain.OrderRepository
	publisher domain.OrderPublisher
}

func NewConfirmReviseOrderHandler(orderRepo domain.OrderRepository, orderPublisher domain.OrderPublisher) ConfirmReviseOrderHandler {
	return ConfirmReviseOrderHandler{
		repo:      orderRepo,
		publisher: orderPublisher,
	}
}

func (h ConfirmReviseOrderHandler) Handle(ctx context.Context, cmd ConfirmReviseOrder) error {
	root, err := h.repo.Update(ctx, cmd.OrderID, &domain.ConfirmReviseOrder{
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return err
	}

	return h.publisher.PublishEntityEvents(ctx, root)
}
