package main

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type orderEventHandlers struct{ app Application }

func newOrderEventHandlers(app Application) orderEventHandlers { return orderEventHandlers{app: app} }

func (h orderEventHandlers) OrderCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*orderapi.OrderCreated)

	err := h.app.Commands.StartCreateOrderSaga.Handle(ctx, commands.StartCreateOrderSaga{
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
