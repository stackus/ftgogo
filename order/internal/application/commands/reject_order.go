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
	counter   domain.Counter
}

func NewRejectOrderHandler(orderRepo domain.OrderRepository, orderPublisher domain.OrderPublisher, ordersRejectedCounter domain.Counter) RejectOrderHandler {
	return RejectOrderHandler{
		repo:      orderRepo,
		publisher: orderPublisher,
		counter:   ordersRejectedCounter,
	}
}

func (h RejectOrderHandler) Handle(ctx context.Context, cmd RejectOrder) error {
	root, err := h.repo.Update(ctx, cmd.OrderID, &domain.RejectOrder{})
	if err != nil {
		return err
	}

	h.counter.Inc()

	return h.publisher.PublishEntityEvents(ctx, root)
}
