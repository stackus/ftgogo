package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
	"serviceapis/orderapi"
)

type StartCreateOrderSaga struct {
	OrderID      string
	ConsumerID   string
	RestaurantID string
	LineItems    []orderapi.LineItem
	OrderTotal   int
}

type StartCreateOrderSagaHandler struct {
	saga    domain.CreateOrderSaga
	counter domain.Counter
}

func NewStartCreateOrderSagaHandler(createOrderSaga domain.CreateOrderSaga, ordersPlaced domain.Counter) StartCreateOrderSagaHandler {
	return StartCreateOrderSagaHandler{
		saga:    createOrderSaga,
		counter: ordersPlaced,
	}
}

func (h StartCreateOrderSagaHandler) Handle(ctx context.Context, cmd StartCreateOrderSaga) error {
	_, err := h.saga.Start(ctx, &domain.CreateOrderSagaData{
		OrderID:      cmd.OrderID,
		ConsumerID:   cmd.ConsumerID,
		RestaurantID: cmd.RestaurantID,
		LineItems:    cmd.LineItems,
		OrderTotal:   cmd.OrderTotal,
	})

	h.counter.Inc()

	return err
}
