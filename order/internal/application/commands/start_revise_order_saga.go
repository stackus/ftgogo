package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type StartReviseOrderSaga struct {
	OrderID           string
	RevisedQuantities map[string]int
}

type StartReviseOrderSagaHandler struct {
	repo domain.OrderRepository
	saga domain.ReviseOrderSaga
}

func NewStartReviseOrderSagaHandler(orderRepo domain.OrderRepository, reviseOrderSaga domain.ReviseOrderSaga) StartReviseOrderSagaHandler {
	return StartReviseOrderSagaHandler{
		repo: orderRepo,
		saga: reviseOrderSaga,
	}
}

func (h StartReviseOrderSagaHandler) Handle(ctx context.Context, cmd StartReviseOrderSaga) (orderapi.OrderState, error) {
	root, err := h.repo.Load(ctx, cmd.OrderID)
	if err != nil {
		return orderapi.UnknownOrderState, err
	}

	order := root.Aggregate().(*domain.Order)

	_, err = h.saga.Start(ctx, &domain.ReviseOrderSagaData{
		OrderID:           cmd.OrderID,
		ConsumerID:        order.ConsumerID,
		RestaurantID:      order.RestaurantID,
		TicketID:          order.TicketID,
		RevisedQuantities: cmd.RevisedQuantities,
	})

	return orderapi.RevisionPending, err
}
