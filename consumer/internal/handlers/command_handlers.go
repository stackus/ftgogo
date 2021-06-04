package handlers

import (
	"context"

	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/application"
	"github.com/stackus/ftgogo/consumer/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/consumerapi"
)

type CommandHandlers struct{ app application.Service }

func NewCommandHandlers(app application.Service) CommandHandlers { return CommandHandlers{app: app} }

func (h CommandHandlers) Mount(subscriber *msg.Subscriber, publisher *msg.Publisher) {
	subscriber.Subscribe(consumerapi.ConsumerServiceCommandChannel, saga.NewCommandDispatcher(publisher).
		Handle(consumerapi.ValidateOrderByConsumer{}, h.ValidateOrderByConsumer))
}

func (h CommandHandlers) ValidateOrderByConsumer(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*consumerapi.ValidateOrderByConsumer)

	err := h.app.Commands.ValidateOrderByConsumer.Handle(ctx, commands.ValidateOrderByConsumer{
		ConsumerID: cmd.ConsumerID,
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return []msg.Reply{msg.WithReply(&consumerapi.ConsumerNotFound{}).Failure()}, nil
		}

		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}
