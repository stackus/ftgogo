package main

import (
	"context"
	"github.com/stackus/edat/msg"
	"github.com/stackus/ftgogo/account/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/consumerapi"
)

type ConsumerEventHandlers struct {
	app Application
}

func NewConsumerEventHandlers(app Application) ConsumerEventHandlers {
	return ConsumerEventHandlers{app: app}
}

func (h ConsumerEventHandlers) ConsumerRegistered(ctx context.Context, evtMsg msg.EntityEvent) error {
	evt := evtMsg.Event().(*consumerapi.ConsumerRegistered)

	return h.app.Commands.CreateAccount.Handle(ctx, commands.CreateAccount{
		ConsumerID: evtMsg.EntityID(),
		Name:       evt.Name,
	})
}
