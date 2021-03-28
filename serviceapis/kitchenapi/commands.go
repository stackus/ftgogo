package kitchenapi

import (
	"github.com/stackus/edat/core"
)

func registerCommands() {
	core.RegisterCommands(
		CreateTicket{}, ConfirmCreateTicket{}, CancelCreateTicket{},
		BeginCancelTicket{}, ConfirmCancelTicket{}, UndoCancelTicket{},
		BeginReviseTicket{}, ConfirmReviseTicket{}, UndoReviseTicket{},
	)
}

type KitchenServiceCommand struct{}

func (KitchenServiceCommand) DestinationChannel() string { return KitchenServiceCommandChannel }

type CreateTicket struct {
	KitchenServiceCommand
	OrderID       string
	RestaurantID  string
	TicketDetails []LineItem
}

func (CreateTicket) CommandName() string { return "kitchenapi.CreateTicket" }

type ConfirmCreateTicket struct {
	KitchenServiceCommand
	TicketID string
}

func (ConfirmCreateTicket) CommandName() string { return "kitchenapi.ConfirmCreateTicket" }

type CancelCreateTicket struct {
	KitchenServiceCommand
	TicketID string
}

func (CancelCreateTicket) CommandName() string { return "kitchenapi.CancelCreateTicket" }

type BeginCancelTicket struct {
	KitchenServiceCommand
	TicketID     string
	RestaurantID string
}

func (BeginCancelTicket) CommandName() string { return "kitchenapi.BeginCancelTicket" }

type ConfirmCancelTicket struct {
	KitchenServiceCommand
	TicketID     string
	RestaurantID string
}

func (ConfirmCancelTicket) CommandName() string { return "kitchenapi.ConfirmCancelTicket" }

type UndoCancelTicket struct {
	KitchenServiceCommand
	TicketID     string
	RestaurantID string
}

func (UndoCancelTicket) CommandName() string { return "kitchenapi.UndoCancelTicket" }

type BeginReviseTicket struct {
	KitchenServiceCommand
	TicketID          string
	RestaurantID      string
	RevisedQuantities map[string]int
}

func (BeginReviseTicket) CommandName() string { return "kitchenapi.BeginReviseTicket" }

type ConfirmReviseTicket struct {
	KitchenServiceCommand
	TicketID          string
	RestaurantID      string
	RevisedQuantities map[string]int
}

func (ConfirmReviseTicket) CommandName() string { return "kitchenapi.ConfirmReviseTicket" }

type UndoReviseTicket struct {
	KitchenServiceCommand
	TicketID     string
	RestaurantID string
}

func (UndoReviseTicket) CommandName() string { return "kitchenapi.UndoReviseTicket" }
