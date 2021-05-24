package main

import (
	"context"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"
	"github.com/stackus/ftgogo/order/internal/application/commands"
	"serviceapis/orderapi"
)

type CommandHandlers struct{ app Application }

func NewCommandHandlers(app Application) CommandHandlers { return CommandHandlers{app: app} }

func (h CommandHandlers) RejectOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*orderapi.RejectOrder)

	err := h.app.Commands.RejectOrder.Handle(ctx, commands.RejectOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ApproveOrder(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*orderapi.ApproveOrder)

	err := h.app.Commands.ApproveOrder.Handle(ctx, commands.ApproveOrder{
		OrderID:  cmd.OrderID,
		TicketID: cmd.TicketID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) BeginCancel(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.BeginCancelOrder)

	err := h.app.Commands.BeginCancelOrder.Handle(ctx, commands.BeginCancelOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) UndoCancel(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.UndoCancelOrder)

	err := h.app.Commands.UndoCancelOrder.Handle(ctx, commands.UndoCancelOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ConfirmCancel(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.ConfirmCancelOrder)

	err := h.app.Commands.ConfirmCancelOrder.Handle(ctx, commands.ConfirmCancelOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) BeginRevise(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.BeginReviseOrder)

	newTotal, err := h.app.Commands.BeginReviseOrder.Handle(ctx, commands.BeginReviseOrder{
		OrderID:           cmd.OrderID,
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithReply(&orderapi.BeginReviseOrderReply{RevisedOrderTotal: newTotal}).Success()}, nil
}

func (h CommandHandlers) UndoRevise(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.UndoReviseOrder)

	err := h.app.Commands.UndoReviseOrder.Handle(ctx, commands.UndoReviseOrder{OrderID: cmd.OrderID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ConfirmRevise(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(orderapi.ConfirmReviseOrder)

	err := h.app.Commands.ConfirmReviseOrder.Handle(ctx, commands.ConfirmReviseOrder{
		OrderID:           cmd.OrderID,
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}
