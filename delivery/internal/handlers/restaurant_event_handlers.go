package handlers

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

type RestaurantEventHandlers struct {
	app application.ServiceApplication
}

func NewRestaurantEventHandlers(app application.ServiceApplication) RestaurantEventHandlers {
	return RestaurantEventHandlers{app: app}
}

func (h RestaurantEventHandlers) Mount(subscriber *msg.Subscriber) {
	subscriber.Subscribe(restaurantapi.RestaurantAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(restaurantapi.RestaurantCreated{}, h.RestaurantCreated))
}

func (h RestaurantEventHandlers) RestaurantCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*restaurantapi.RestaurantCreated)

	return h.app.CreateRestaurant(ctx, commands.CreateRestaurant{
		RestaurantID: evtMsg.EntityID(),
		Name:         evt.Name,
		Address: &commonapi.Address{
			Street1: evt.Address.Street1,
			Street2: evt.Address.Street2,
			City:    evt.Address.City,
			State:   evt.Address.State,
			Zip:     evt.Address.Zip,
		},
	})
}
