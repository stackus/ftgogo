package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
)

type ApproveOrder struct {
	OrderID  string
	TicketID string
}

type ApproveOrderHandler struct {
	repo    ports.OrderRepository
	counter ports.Counter
}

func NewApproveOrderHandler(repo ports.OrderRepository, counter ports.Counter) ApproveOrderHandler {
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
