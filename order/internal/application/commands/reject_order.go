package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type RejectOrder struct {
	OrderID string
}

type RejectOrderHandler struct {
	repo    domain.OrderRepository
	counter domain.Counter
}

func NewRejectOrderHandler(repo domain.OrderRepository, counter domain.Counter) RejectOrderHandler {
	return RejectOrderHandler{
		repo:    repo,
		counter: counter,
	}
}

func (h RejectOrderHandler) Handle(ctx context.Context, cmd RejectOrder) error {
	_, err := h.repo.Update(ctx, cmd.OrderID, &domain.RejectOrder{})
	if err != nil {
		return err
	}

	h.counter.Inc()

	return nil
}
