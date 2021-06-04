package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type ApproveOrder struct {
	OrderID  string
	TicketID string
}

type ApproveOrderHandler struct {
	repo    adapters.OrderRepository
	counter adapters.Counter
}

func NewApproveOrderHandler(repo adapters.OrderRepository, counter adapters.Counter) ApproveOrderHandler {
	return ApproveOrderHandler{
		repo:    repo,
		counter: counter,
	}
}

func (h ApproveOrderHandler) Handle(ctx context.Context, cmd ApproveOrder) error {
	_, err := h.repo.Update(ctx, cmd.OrderID, &domain.ApproveOrder{TicketID: cmd.TicketID})
	if err != nil {
		return err
	}

	h.counter.Inc()

	return nil
}
