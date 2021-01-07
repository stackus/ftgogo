package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type RejectOrder struct {
	OrderID string
}

type RejectOrderHandler struct {
	repo      domain.OrderRepository
	publisher domain.OrderPublisher
}

func NewRejectOrderHandler(orderRepo domain.OrderRepository, orderPublisher domain.OrderPublisher) RejectOrderHandler {
	return RejectOrderHandler{
		repo:      orderRepo,
		publisher: orderPublisher,
	}
}

func (h RejectOrderHandler) Handle(ctx context.Context, cmd RejectOrder) error {
	root, err := h.repo.Update(ctx, cmd.OrderID, &domain.RejectOrder{})
	if err != nil {
		return err
	}

	return h.publisher.PublishEntityEvents(ctx, root)
}
