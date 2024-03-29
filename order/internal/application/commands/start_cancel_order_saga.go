package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type StartCancelOrderSaga struct {
	OrderID string
}

type StartCancelOrderSagaHandler struct {
	repo ports.OrderRepository
	saga ports.CancelOrderSaga
}

func NewStartCancelOrderSagaHandler(orderRepo ports.OrderRepository, cancelOrderSaga ports.CancelOrderSaga) StartCancelOrderSagaHandler {
	return StartCancelOrderSagaHandler{
		repo: orderRepo,
		saga: cancelOrderSaga,
	}
}

func (h StartCancelOrderSagaHandler) Handle(ctx context.Context, cmd StartCancelOrderSaga) (orderapi.OrderState, error) {
	order, err := h.repo.Load(ctx, cmd.OrderID)
	if err != nil {
		return orderapi.UnknownOrderState, err
	}

	_, err = h.saga.Start(ctx, &domain.CancelOrderSagaData{
		OrderID:      cmd.OrderID,
		ConsumerID:   order.ConsumerID,
		RestaurantID: order.RestaurantID,
		TicketID:     order.TicketID,
		OrderTotal:   order.OrderTotal(),
	})

	return orderapi.CancelPending, err
}
