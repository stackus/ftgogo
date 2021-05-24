package domain

import (
	"time"

	"github.com/stackus/edat/core"

	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
)

func registerTicketCommands() {
	core.RegisterCommands(
		CreateTicket{}, ConfirmCreateTicket{}, CancelCreateTicket{},
		CancelTicket{}, ConfirmCancelTicket{}, UndoCancelTicket{},
		ReviseTicket{}, ConfirmReviseTicket{}, UndoReviseTicket{},
		AcceptTicket{},
	)
}

type CreateTicket struct {
	OrderID      string
	RestaurantID string
	LineItems    []kitchenapi.LineItem
}

func (CreateTicket) CommandName() string { return "kitchenservice.CreateTicket" }

type ConfirmCreateTicket struct{}

func (ConfirmCreateTicket) CommandName() string { return "kitchenservice.ConfirmCreateTicket" }

type CancelCreateTicket struct{}

func (CancelCreateTicket) CommandName() string { return "kitchenservice.CancelCreateTicket" }

type CancelTicket struct{}

func (CancelTicket) CommandName() string { return "kitchenservice.CancelTicket" }

type ConfirmCancelTicket struct{}

func (ConfirmCancelTicket) CommandName() string { return "kitchenservice.ConfirmCancelTicket" }

type UndoCancelTicket struct{}

func (UndoCancelTicket) CommandName() string { return "kitchenservice.UndoCancelTicket" }

type ReviseTicket struct{}

func (ReviseTicket) CommandName() string { return "kitchenservice.ReviseTicket" }

type ConfirmReviseTicket struct {
	RevisedQuantities map[string]int
}

func (ConfirmReviseTicket) CommandName() string { return "kitchenservice.ConfirmReviseTicket" }

type UndoReviseTicket struct{}

func (UndoReviseTicket) CommandName() string { return "kitchenservice.UndoReviseTicket" }

type AcceptTicket struct {
	ReadyBy time.Time
}

func (AcceptTicket) CommandName() string { return "kitchenservice.AcceptTicket" }
