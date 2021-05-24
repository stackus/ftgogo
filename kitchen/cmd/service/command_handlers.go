package main

import (
	"context"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"
	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
)

type CommandHandlers struct{ app Application }

func NewCommandHandlers(app Application) CommandHandlers { return CommandHandlers{app: app} }

func (h CommandHandlers) CreateTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.CreateTicket)

	ticketID, err := h.app.Commands.CreateTicket.Handle(ctx, commands.CreateTicket{
		OrderID:      cmd.OrderID,
		RestaurantID: cmd.RestaurantID,
		LineItems:    cmd.TicketDetails,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithReply(&kitchenapi.CreateTicketReply{TicketID: ticketID}).Success()}, nil
}

func (h CommandHandlers) ConfirmCreateTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.ConfirmCreateTicket)

	err := h.app.Commands.ConfirmCreateTicket.Handle(ctx, commands.ConfirmCreateTicket{TicketID: cmd.TicketID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) CancelCreateTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.CancelCreateTicket)

	err := h.app.Commands.CancelCreateTicket.Handle(ctx, commands.CancelCreateTicket{TicketID: cmd.TicketID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) BeginCancelTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.BeginCancelTicket)

	err := h.app.Commands.BeginCancelTicket.Handle(ctx, commands.BeginCancelTicket{
		TicketID:     cmd.TicketID,
		RestaurantID: cmd.RestaurantID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ConfirmCancelTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.ConfirmCancelTicket)

	err := h.app.Commands.ConfirmCancelTicket.Handle(ctx, commands.ConfirmCancelTicket{
		TicketID:     cmd.TicketID,
		RestaurantID: cmd.RestaurantID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) UndoCancelTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.UndoCancelTicket)

	err := h.app.Commands.UndoCancelTicket.Handle(ctx, commands.UndoCancelTicket{
		TicketID:     cmd.TicketID,
		RestaurantID: cmd.RestaurantID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) BeginReviseTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.BeginReviseTicket)

	err := h.app.Commands.BeginReviseTicket.Handle(ctx, commands.BeginReviseTicket{
		TicketID:          cmd.TicketID,
		RestaurantID:      cmd.RestaurantID,
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) ConfirmReviseTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.ConfirmReviseTicket)

	err := h.app.Commands.ConfirmReviseTicket.Handle(ctx, commands.ConfirmReviseTicket{
		TicketID:          cmd.TicketID,
		RestaurantID:      cmd.RestaurantID,
		RevisedQuantities: cmd.RevisedQuantities,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) UndoReviseTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.UndoReviseTicket)

	err := h.app.Commands.UndoReviseTicket.Handle(ctx, commands.UndoReviseTicket{
		TicketID:     cmd.TicketID,
		RestaurantID: cmd.RestaurantID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}
