package handlers

import (
	"context"

	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/accounting/internal/application"
	"github.com/stackus/ftgogo/accounting/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/consumerapi"
)

type ConsumerEventHandlers struct {
	app application.Service
}

func NewConsumerEventHandlers(app application.Service) ConsumerEventHandlers {
	return ConsumerEventHandlers{app: app}
}

func (h ConsumerEventHandlers) Mount(subscriber *msg.Subscriber) {
	subscriber.Subscribe(consumerapi.ConsumerAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(consumerapi.ConsumerRegistered{}, h.ConsumerRegistered))
}

func (h ConsumerEventHandlers) ConsumerRegistered(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*consumerapi.ConsumerRegistered)

	return h.app.Commands.CreateAccount.Handle(ctx, commands.CreateAccount{
		ConsumerID: evtMsg.EntityID(),
		Name:       evt.Name,
	})
}
