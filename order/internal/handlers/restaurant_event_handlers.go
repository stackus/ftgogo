package handlers

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/order/internal/application"
	"github.com/stackus/ftgogo/order/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

type RestaurantEventHandlers struct{ app application.Service }

func NewRestaurantEventHandlers(app application.Service) RestaurantEventHandlers {
	return RestaurantEventHandlers{app: app}
}

func (h RestaurantEventHandlers) Mount(subscriber *msg.Subscriber) {
	subscriber.Subscribe(restaurantapi.RestaurantAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(restaurantapi.RestaurantCreated{}, h.RestaurantCreated).
		Handle(restaurantapi.RestaurantMenuRevised{}, h.RestaurantMenuRevised))
}

func (h RestaurantEventHandlers) RestaurantCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*restaurantapi.RestaurantCreated)

	return h.app.Commands.CreateRestaurant.Handle(ctx, commands.CreateRestaurant{
		RestaurantID: evtMsg.EntityID(),
		Name:         evt.Name,
		Menu:         evt.Menu,
	})
}

func (h RestaurantEventHandlers) RestaurantMenuRevised(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*restaurantapi.RestaurantMenuRevised)

	return h.app.Commands.ReviseRestaurantMenu.Handle(ctx, commands.ReviseRestaurantMenu{
		RestaurantID: evtMsg.EntityID(),
		Menu:         evt.Menu,
	})
}
