package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type StartCreateOrderSaga struct {
	OrderID      string
	ConsumerID   string
	RestaurantID string
	LineItems    []orderapi.LineItem
	OrderTotal   int
}

type StartCreateOrderSagaHandler struct {
	saga    ports.CreateOrderSaga
	counter ports.Counter
}

func NewStartCreateOrderSagaHandler(createOrderSaga ports.CreateOrderSaga, ordersPlaced ports.Counter) StartCreateOrderSagaHandler {
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
