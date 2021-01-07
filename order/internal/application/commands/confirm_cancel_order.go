package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type ConfirmCancelOrder struct {
	OrderID string
}

type ConfirmCancelOrderHandler struct {
	repo      domain.OrderRepository
	publisher domain.OrderPublisher
}

func NewConfirmCancelOrderHandler(orderRepo domain.OrderRepository, orderPublisher domain.OrderPublisher) ConfirmCancelOrderHandler {
	return ConfirmCancelOrderHandler{
		repo:      orderRepo,
		publisher: orderPublisher,
	}
}

func (h ConfirmCancelOrderHandler) Handle(ctx context.Context, cmd ConfirmCancelOrder) error {
	root, err := h.repo.Update(ctx, cmd.OrderID, &domain.ConfirmCancelOrder{})
	if err != nil {
		return err
	}

	return h.publisher.PublishEntityEvents(ctx, root)
}
