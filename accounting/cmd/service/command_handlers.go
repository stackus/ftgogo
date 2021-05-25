package main

import (
	"context"
	"errors"

	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/account/internal/application/commands"
	"github.com/stackus/ftgogo/account/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/accountingapi"
)

type CommandHandlers struct {
	app Application
}

func NewCommandHandlers(app Application) CommandHandlers {
	return CommandHandlers{app: app}
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
