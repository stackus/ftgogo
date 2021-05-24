package main

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/deliveryapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

type restaurantEventHandlers struct{ app Application }

func newRestaurantEventHandlers(app Application) restaurantEventHandlers {
	return restaurantEventHandlers{app: app}
}

func (h restaurantEventHandlers) RestaurantCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*restaurantapi.RestaurantCreated)

	return h.app.Commands.CreateRestaurant.Handle(ctx, commands.CreateRestaurant{
		RestaurantID: evtMsg.EntityID(),
		Name:         evt.Name,
		Address: deliveryapi.Address{
			Street1: evt.Address.Street1,
			Street2: evt.Address.Street2,
			City:    evt.Address.City,
			State:   evt.Address.State,
			Zip:     evt.Address.Zip,
		},
	})
}
