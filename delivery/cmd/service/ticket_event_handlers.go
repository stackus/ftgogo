package main

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
)

type ticketEventHandlers struct{ app Application }

func newTicketEventHandlers(app Application) ticketEventHandlers {
	return ticketEventHandlers{app: app}
}

func (h ticketEventHandlers) TicketAccepted(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*kitchenapi.TicketAccepted)

	return h.app.Commands.ScheduleDelivery.Handle(ctx, commands.ScheduleDelivery{
		OrderID: evt.OrderID,
		ReadyBy: evt.ReadyBy,
	})
}

func (h ticketEventHandlers) TicketCancelled(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*kitchenapi.TicketCancelled)

	return h.app.Commands.CancelDelivery.Handle(ctx, commands.CancelDelivery{OrderID: evt.OrderID})
}
