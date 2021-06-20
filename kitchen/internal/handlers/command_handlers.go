package handlers

import (
	"context"

	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/kitchen/internal/application"
	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
)

type CommandHandlers struct {
	app application.ServiceApplication
}

func NewCommandHandlers(app application.ServiceApplication) CommandHandlers {
	return CommandHandlers{app: app}
}

func (h CommandHandlers) Mount(subscriber *msg.Subscriber, publisher *msg.Publisher) {
	subscriber.Subscribe(kitchenapi.KitchenServiceCommandChannel, saga.NewCommandDispatcher(publisher).
		Handle(kitchenapi.CreateTicket{}, h.CreateTicket).
		Handle(kitchenapi.ConfirmCreateTicket{}, h.ConfirmCreateTicket).
		Handle(kitchenapi.CancelCreateTicket{}, h.CancelCreateTicket).
		Handle(kitchenapi.BeginCancelTicket{}, h.BeginCancelTicket).
		Handle(kitchenapi.ConfirmCancelTicket{}, h.ConfirmCancelTicket).
		Handle(kitchenapi.UndoCancelTicket{}, h.UndoCancelTicket).
		Handle(kitchenapi.BeginReviseTicket{}, h.BeginReviseTicket).
		Handle(kitchenapi.ConfirmReviseTicket{}, h.ConfirmReviseTicket).
		Handle(kitchenapi.UndoReviseTicket{}, h.UndoReviseTicket))
}

func (h CommandHandlers) CreateTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.CreateTicket)

	ticketID, err := h.app.CreateTicket(ctx, commands.CreateTicket{
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

	err := h.app.ConfirmCreateTicket(ctx, commands.ConfirmCreateTicket{TicketID: cmd.TicketID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) CancelCreateTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.CancelCreateTicket)

	err := h.app.CancelCreateTicket(ctx, commands.CancelCreateTicket{TicketID: cmd.TicketID})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}

func (h CommandHandlers) BeginCancelTicket(ctx context.Context, cmdMsg saga.Command) ([]msg.Reply, error) {
	cmd := cmdMsg.Command().(*kitchenapi.BeginCancelTicket)

	err := h.app.BeginCancelTicket(ctx, commands.BeginCancelTicket{
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

	err := h.app.ConfirmCancelTicket(ctx, commands.ConfirmCancelTicket{
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

	err := h.app.UndoCancelTicket(ctx, commands.UndoCancelTicket{
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

	err := h.app.BeginReviseTicket(ctx, commands.BeginReviseTicket{
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

	err := h.app.ConfirmReviseTicket(ctx, commands.ConfirmReviseTicket{
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

	err := h.app.UndoReviseTicket(ctx, commands.UndoReviseTicket{
		TicketID:     cmd.TicketID,
		RestaurantID: cmd.RestaurantID,
	})
	if err != nil {
		return []msg.Reply{msg.WithFailure()}, nil
	}

	return []msg.Reply{msg.WithSuccess()}, nil
}
