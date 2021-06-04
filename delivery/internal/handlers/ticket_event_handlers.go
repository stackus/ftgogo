package handlers

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
)

type TicketEventHandlers struct{ app application.Service }

func NewTicketEventHandlers(app application.Service) TicketEventHandlers {
	return TicketEventHandlers{app: app}
}

func (h TicketEventHandlers) Mount(subscriber *msg.Subscriber) {
	subscriber.Subscribe(kitchenapi.TicketAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(kitchenapi.TicketAccepted{}, h.TicketAccepted).
		Handle(kitchenapi.TicketCancelled{}, h.TicketCancelled))
}

func (h TicketEventHandlers) TicketAccepted(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*kitchenapi.TicketAccepted)

	return h.app.Commands.ScheduleDelivery.Handle(ctx, commands.ScheduleDelivery{
		OrderID: evt.OrderID,
		ReadyBy: evt.ReadyBy,
	})
}

func (h TicketEventHandlers) TicketCancelled(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*kitchenapi.TicketCancelled)

	return h.app.Commands.CancelDelivery.Handle(ctx, commands.CancelDelivery{OrderID: evt.OrderID})
}
