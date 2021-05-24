package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type StartCancelOrderSaga struct {
	OrderID string
}

type StartCancelOrderSagaHandler struct {
	repo domain.OrderRepository
	saga domain.CancelOrderSaga
}

func NewStartCancelOrderSagaHandler(orderRepo domain.OrderRepository, cancelOrderSaga domain.CancelOrderSaga) StartCancelOrderSagaHandler {
	return StartCancelOrderSagaHandler{
		repo: orderRepo,
		saga: cancelOrderSaga,
	}
}

func (h StartCancelOrderSagaHandler) Handle(ctx context.Context, cmd StartCancelOrderSaga) (orderapi.OrderState, error) {
	root, err := h.repo.Load(ctx, cmd.OrderID)
	if err != nil {
		return orderapi.UnknownOrderState, err
	}

	order := root.Aggregate().(*domain.Order)

	_, err = h.saga.Start(ctx, &domain.CancelOrderSagaData{
		OrderID:      cmd.OrderID,
		ConsumerID:   order.ConsumerID,
		RestaurantID: order.RestaurantID,
		TicketID:     order.TicketID,
		OrderTotal:   order.OrderTotal(),
	})

	return orderapi.CancelPending, err
}
