package main

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/order-history/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type OrderEventHandlers struct{ app Application }

func NewOrderEventHandlers(app Application) OrderEventHandlers { return OrderEventHandlers{app: app} }

func (h OrderEventHandlers) OrderCreated(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*orderapi.OrderCreated)

	return h.app.Commands.CreateOrderHistory.Handle(ctx, commands.CreateOrderHistory{
		OrderID:        evtMsg.EntityID(),
		ConsumerID:     evt.ConsumerID,
		RestaurantID:   evt.RestaurantID,
		RestaurantName: evt.RestaurantName,
		LineItems:      evt.LineItems,
		OrderTotal:     evt.OrderTotal,
	})
}

func (h OrderEventHandlers) OrderApproved(ctx context.Context, evtMsg msg.EntityEvent) error {
	return h.app.Commands.UpdateOrderStatus.Handle(ctx, commands.UpdateOrderStatus{
		OrderID: evtMsg.EntityID(),
		Status:  orderapi.Approved,
	})
}

func (h OrderEventHandlers) OrderCancelled(ctx context.Context, evtMsg msg.EntityEvent) error {
	return h.app.Commands.UpdateOrderStatus.Handle(ctx, commands.UpdateOrderStatus{
		OrderID: evtMsg.EntityID(),
		Status:  orderapi.Cancelled,
	})
}

func (h OrderEventHandlers) OrderRejected(ctx context.Context, evtMsg msg.EntityEvent) error {
	return h.app.Commands.UpdateOrderStatus.Handle(ctx, commands.UpdateOrderStatus{
		OrderID: evtMsg.EntityID(),
		Status:  orderapi.Rejected,
	})
}
