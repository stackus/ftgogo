package handlers

import (
	"context"
	"errors"

	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/accounting/internal/application"
	"github.com/stackus/ftgogo/accounting/internal/application/commands"
	"github.com/stackus/ftgogo/accounting/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/accountingapi"
)

type CommandHandlers struct {
	app application.Service
}

func NewCommandHandlers(app application.Service) CommandHandlers {
	return CommandHandlers{app: app}
}
func (h CommandHandlers) Mount(subscriber *msg.Subscriber, publisher *msg.Publisher) {
	subscriber.Subscribe(accountingapi.AccountingServiceCommandChannel, saga.NewCommandDispatcher(publisher).
		Handle(accountingapi.AuthorizeOrder{}, h.AuthorizeOrder).
		Handle(accountingapi.ReverseAuthorizeOrder{}, h.ReverseAuthorizeOrder).
		Handle(accountingapi.ReviseAuthorizeOrder{}, h.ReviseAuthorizeOrder))
}

func (h CommandHandlers) AuthorizeOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*accountingapi.AuthorizeOrder)

	err := h.app.Commands.AuthorizeOrder.Handle(ctx, commands.AuthorizeOrder{
		ConsumerID: cmd.ConsumerID,
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAccountDisabled) {
			return []msg.Reply{msg.WithReply(&accountingapi.AccountDisabled{}).Failure()}, nil
		}
		return nil, err
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ReverseAuthorizeOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*accountingapi.ReverseAuthorizeOrder)

	err := h.app.Commands.ReverseAuthorizeOrder.Handle(ctx, commands.ReverseAuthorizeOrder{
		ConsumerID: cmd.ConsumerID,
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAccountDisabled) {
			return []msg.Reply{msg.WithReply(&accountingapi.AccountDisabled{}).Failure()}, nil
		}
		return nil, err
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ReviseAuthorizeOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*accountingapi.ReviseAuthorizeOrder)

	err := h.app.Commands.ReviseAuthorizeOrder.Handle(ctx, commands.ReviseAuthorizeOrder{
		ConsumerID: cmd.ConsumerID,
		OrderID:    cmd.OrderID,
		OrderTotal: cmd.OrderTotal,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAccountDisabled) {
			return []msg.Reply{msg.WithReply(&accountingapi.AccountDisabled{}).Failure()}, nil
		}
		return nil, err
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}
