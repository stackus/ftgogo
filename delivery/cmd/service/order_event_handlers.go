package main

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/deliveryapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type orderEventHandlers struct{ app Application }

func newOrderEventHandlers(app Application) orderEventHandlers { return orderEventHandlers{app: app} }

func (h orderEventHandlers) OrderCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*orderapi.OrderCreated)

	return h.app.Commands.CreateDelivery.Handle(ctx, commands.CreateDelivery{
		OrderID:      evtMsg.EntityID(),
		RestaurantID: evt.RestaurantID,
		DeliveryAddress: deliveryapi.Address{
			Street1: evt.DeliverTo.Street1,
			Street2: evt.DeliverTo.Street2,
			City:    evt.DeliverTo.City,
			State:   evt.DeliverTo.State,
			Zip:     evt.DeliverTo.Zip,
		},
	})
}
