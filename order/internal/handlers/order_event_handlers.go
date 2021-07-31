package handlers

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/order/internal/application"
	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type OrderEventHandlers struct {
	app application.ServiceApplication
}

func NewOrderEventHandlers(app application.ServiceApplication) OrderEventHandlers {
	return OrderEventHandlers{app: app}
}

func (h OrderEventHandlers) Mount(subscriber *msg.Subscriber) {
	subscriber.Subscribe(orderapi.OrderAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(orderapi.OrderCreated{}, h.OrderCreated))
}

func (h OrderEventHandlers) OrderCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*orderapi.OrderCreated)

	err := h.app.StartCreateOrderSaga(ctx, commands.StartCreateOrderSaga{
		OrderID:      evtMsg.EntityID(),
		ConsumerID:   evt.ConsumerID,
		RestaurantID: evt.RestaurantID,
		LineItems:    evt.LineItems,
		OrderTotal:   evt.OrderTotal,
	})
	if err != nil {
		return err
	}

	return nil
}
