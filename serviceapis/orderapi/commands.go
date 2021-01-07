package orderapi

import (
	"github.com/stackus/edat/core"
)

func registerCommands() {
	core.RegisterCommands(
		RejectOrder{}, ApproveOrder{},
		BeginCancelOrder{}, UndoCancelOrder{}, ConfirmCancelOrder{},
		BeginReviseOrder{}, UndoReviseOrder{}, ConfirmReviseOrder{},
	)
}

type OrderServiceCommand struct{}

func (OrderServiceCommand) DestinationChannel() string { return OrderServiceCommandChannel }

type RejectOrder struct {
	OrderServiceCommand
	OrderID string
}

func (RejectOrder) CommandName() string { return "orderapi.RejectOrder" }

type ApproveOrder struct {
	OrderServiceCommand
	OrderID  string
	TicketID string
}

func (ApproveOrder) CommandName() string { return "orderapi.ApproveOrder" }

// Cancel Order Saga

type CancelOrderSagaCommand struct{}

func (CancelOrderSagaCommand) DestinationChannel() string { return CancelOrderSagaChannel }

type BeginCancelOrder struct {
	CancelOrderSagaCommand
	OrderID string
}

func (BeginCancelOrder) CommandName() string { return "orderapi.BeginCancelOrder" }

type UndoCancelOrder struct {
	CancelOrderSagaCommand
	OrderID string
}

func (UndoCancelOrder) CommandName() string { return "orderapi.UndoCancelOrder" }

type ConfirmCancelOrder struct {
	CancelOrderSagaCommand
	OrderID string
}

func (ConfirmCancelOrder) CommandName() string { return "orderapi.ConfirmCancelOrder" }

// Revise Order Saga

type ReviseOrderSagaCommand struct{}

func (ReviseOrderSagaCommand) DestinationChannel() string { return ReviseOrderSagaChannel }

type BeginReviseOrder struct {
	ReviseOrderSagaCommand
	OrderID           string
	RevisedQuantities map[string]int
}

func (BeginReviseOrder) CommandName() string { return "orderapi.BeginReviseOrder" }

type UndoReviseOrder struct {
	ReviseOrderSagaCommand
	OrderID string
}

func (UndoReviseOrder) CommandName() string { return "orderapi.UndoReviseOrder" }

type ConfirmReviseOrder struct {
	ReviseOrderSagaCommand
	OrderID           string
	RevisedQuantities map[string]int
}

func (ConfirmReviseOrder) CommandName() string { return "orderapi.ConfirmReviseOrder" }
