package handlers

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/deliveryapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type OrderEventHandlers struct{ app application.Service }

func NewOrderEventHandlers(app application.Service) OrderEventHandlers {
	return OrderEventHandlers{app: app}
}

func (h OrderEventHandlers) Mount(subscriber *msg.Subscriber) {
	subscriber.Subscribe(orderapi.OrderAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(orderapi.OrderCreated{}, h.OrderCreated))
}

func (h OrderEventHandlers) OrderCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
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