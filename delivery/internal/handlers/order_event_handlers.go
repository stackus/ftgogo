package handlers

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
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

	return h.app.CreateDelivery(ctx, commands.CreateDelivery{
		OrderID:      evtMsg.EntityID(),
		RestaurantID: evt.RestaurantID,
		DeliveryAddress: &commonapi.Address{
			Street1: evt.DeliverTo.Street1,
			Street2: evt.DeliverTo.Street2,
			City:    evt.DeliverTo.City,
			State:   evt.DeliverTo.State,
			Zip:     evt.DeliverTo.Zip,
		},
	})
}
