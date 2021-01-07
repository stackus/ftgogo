package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type ApproveOrder struct {
	OrderID  string
	TicketID string
}

type ApproveOrderHandler struct {
	repo      domain.OrderRepository
	publisher domain.OrderPublisher
}

func NewApproveOrderHandler(orderRepo domain.OrderRepository, orderPublisher domain.OrderPublisher) ApproveOrderHandler {
	return ApproveOrderHandler{
		repo:      orderRepo,
		publisher: orderPublisher,
	}
}

func (h ApproveOrderHandler) Handle(ctx context.Context, cmd ApproveOrder) error {
	root, err := h.repo.Update(ctx, cmd.OrderID, &domain.ApproveOrder{TicketID: cmd.TicketID})
	if err != nil {
		return err
	}

	return h.publisher.PublishEntityEvents(ctx, root)
}
