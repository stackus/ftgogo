package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type StartReviseOrderSaga struct {
	OrderID           string
	RevisedQuantities map[string]int
}

type StartReviseOrderSagaHandler struct {
	repo adapters.OrderRepository
	saga adapters.ReviseOrderSaga
}

func NewStartReviseOrderSagaHandler(orderRepo adapters.OrderRepository, reviseOrderSaga adapters.ReviseOrderSaga) StartReviseOrderSagaHandler {
	return StartReviseOrderSagaHandler{
		repo: orderRepo,
		saga: reviseOrderSaga,
	}
}

func (h StartReviseOrderSagaHandler) Handle(ctx context.Context, cmd StartReviseOrderSaga) (orderapi.OrderState, error) {
	order, err := h.repo.Load(ctx, cmd.OrderID)
	if err != nil {
		return orderapi.UnknownOrderState, err
	}

	_, err = h.saga.Start(ctx, &domain.ReviseOrderSagaData{
		OrderID:           cmd.OrderID,
		ConsumerID:        order.ConsumerID,
		RestaurantID:      order.RestaurantID,
		TicketID:          order.TicketID,
		RevisedQuantities: cmd.RevisedQuantities,
	})

	return orderapi.RevisionPending, err
}
