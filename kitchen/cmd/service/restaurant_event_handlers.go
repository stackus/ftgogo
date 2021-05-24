package main

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
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
		Menu:         evt.Menu,
	})
}

func (h restaurantEventHandlers) RestaurantMenuRevised(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*restaurantapi.RestaurantMenuRevised)

	return h.app.Commands.ReviseRestaurantMenu.Handle(ctx, commands.ReviseRestaurantMenu{
		RestaurantID: evtMsg.EntityID(),
		Menu:         evt.Menu,
	})
}
