package domain

import (
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"time"

	"github.com/stackus/edat/core"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

func registerOrderCommands() {
	core.RegisterCommands(
		CreateOrder{}, RejectOrder{}, ApproveOrder{},
		BeginCancelOrder{}, UndoCancelOrder{}, ConfirmCancelOrder{},
		BeginReviseOrder{}, UndoReviseOrder{}, ConfirmReviseOrder{},
	)
}

// CreateOrder order command
type CreateOrder struct {
	ConsumerID     string
	RestaurantID   string
	RestaurantName string
	LineItems      []orderapi.LineItem
	OrderTotal     int
	DeliverAt      time.Time
	DeliverTo      commonapi.Address
}

// CommandName command method
func (CreateOrder) CommandName() string { return "orderservice.CreateOrder" }

// RejectOrder order command
type RejectOrder struct{}

// CommandName command method
func (RejectOrder) CommandName() string { return "orderservice.RejectOrder" }

// ApproveOrder order Command
type ApproveOrder struct {
	TicketID string
}

// CommandName command method
func (ApproveOrder) CommandName() string { return "orderservice.ApproveOrder" }

// Cancel

// BeginCancelOrder order command
type BeginCancelOrder struct{}

// CommandName command method
func (BeginCancelOrder) CommandName() string { return "orderservice.BeginCancelOrder" }

// UndoCancelOrder order command
type UndoCancelOrder struct{}

// CommandName command method
func (UndoCancelOrder) CommandName() string { return "orderservice.UndoCancelOrder" }

// ConfirmCancelOrder order command
type ConfirmCancelOrder struct{}

// CommandName command method
func (ConfirmCancelOrder) CommandName() string { return "orderservice.ConfirmCancelOrder" }

// Revise

// BeginReviseOrder order command
type BeginReviseOrder struct {
	RevisedQuantities map[string]int
}

// CommandName command method
func (BeginReviseOrder) CommandName() string { return "orderservice.BeginReviseOrder" }

// UndoReviseOrder order command
type UndoReviseOrder struct{}

// CommandName command method
func (UndoReviseOrder) CommandName() string { return "orderservice.UndoReviseOrder" }

// ConfirmReviseOrder order command
type ConfirmReviseOrder struct {
	RevisedQuantities map[string]int
}

// CommandName command method
func (ConfirmReviseOrder) CommandName() string { return "orderservice.ConfirmReviseOrder" }
