package domain

import (
	"fmt"
	"time"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"serviceapis/commonapi"
	"serviceapis/orderapi"
	"shared-go/errs"
)

var (
	ErrOrderUnhandledCommand  = errs.NewError("unhandled command in order aggregate", errs.ErrServerError)
	ErrOrderUnhandledEvent    = errs.NewError("unhandled event in order aggregate", errs.ErrServerError)
	ErrOrderUnhandledSnapshot = errs.NewError("unhandled snapshot in order aggregate", errs.ErrServerError)

	ErrOrderInvalidState  = errs.NewError("order state does not allow action", errs.ErrConflict)
	ErrOrderMinimumNotMet = errs.NewError("order total does not meet the minimum", errs.ErrUnprocessableEntity)
)

const orderMinimum = 0

type Order struct {
	es.AggregateBase
	ConsumerID   string
	RestaurantID string
	TicketID     string
	LineItems    []orderapi.LineItem
	State        orderapi.OrderState
	DeliverAt    time.Time
	DeliverTo    commonapi.Address
}

var _ es.Aggregate = (*Order)(nil)

func NewOrder() es.Aggregate {
	return &Order{}
}

// EntityName aggregate method
func (Order) EntityName() string {
	return "orderservice.Order"
}

// ProcessCommand aggregate method
func (o *Order) ProcessCommand(command core.Command) error {
	switch cmd := command.(type) {
	case *CreateOrder:
		if o.State != orderapi.UnknownOrderState {
			return ErrOrderInvalidState
		}

		o.AddEvent(&orderapi.OrderCreated{
			ConsumerID:     cmd.ConsumerID,
			RestaurantID:   cmd.RestaurantID,
			RestaurantName: cmd.RestaurantName,
			LineItems:      cmd.LineItems,
			OrderTotal:     cmd.OrderTotal,
			DeliverAt:      cmd.DeliverAt,
			DeliverTo:      cmd.DeliverTo,
		})

	case *ApproveOrder:
		if o.State != orderapi.ApprovalPending {
			return ErrOrderInvalidState
		}
		o.AddEvent(&orderapi.OrderApproved{TicketID: cmd.TicketID})

	case *RejectOrder:
		if o.State != orderapi.ApprovalPending {
			return ErrOrderInvalidState
		}
		o.AddEvent(&orderapi.OrderRejected{})

	case *BeginCancelOrder:
		if o.State != orderapi.Approved {
			return ErrOrderInvalidState
		}
		o.AddEvent(&OrderCancelling{})

	case *UndoCancelOrder:
		if o.State != orderapi.CancelPending {
			return ErrOrderInvalidState
		}
		o.AddEvent(&OrderCancellingUndone{})

	case *ConfirmCancelOrder:
		if o.State != orderapi.CancelPending {
			return ErrOrderInvalidState
		}
		o.AddEvent(&orderapi.OrderCancelled{})

	case *BeginReviseOrder:
		if o.State != orderapi.Approved {
			return ErrOrderInvalidState
		}

		currentTotal := o.OrderTotal()
		newTotal := o.RevisedOrderTotal(currentTotal, cmd.RevisedQuantities)

		if newTotal < orderMinimum {
			return ErrOrderMinimumNotMet
		}

		o.AddEvent(&orderapi.OrderProposedRevision{
			CurrentOrderTotal: currentTotal,
			NewOrderTotal:     newTotal,
			Revisions:         cmd.RevisedQuantities,
		})

	case *UndoReviseOrder:
		if o.State != orderapi.RevisionPending {
			return ErrOrderInvalidState
		}
		o.AddEvent(&OrderRevisingUndone{})

	case *ConfirmReviseOrder:
		if o.State != orderapi.RevisionPending {
			return ErrOrderInvalidState
		}

		currentTotal := o.OrderTotal()
		newTotal := o.RevisedOrderTotal(currentTotal, cmd.RevisedQuantities)

		o.AddEvent(&orderapi.OrderRevised{
			CurrentOrderTotal: currentTotal,
			NewOrderTotal:     newTotal,
			RevisedQuantities: cmd.RevisedQuantities,
		})

	default:
		return errs.NewError(fmt.Sprintf("unhandled command `%T`", command), ErrOrderUnhandledCommand)
	}

	return nil
}

// ApplyEvent aggregate method
func (o *Order) ApplyEvent(event core.Event) error {
	switch evt := event.(type) {
	case *orderapi.OrderCreated:
		o.ConsumerID = evt.ConsumerID
		o.RestaurantID = evt.RestaurantID
		o.LineItems = evt.LineItems
		o.DeliverAt = evt.DeliverAt
		o.DeliverTo = evt.DeliverTo
		o.State = orderapi.ApprovalPending

	case *orderapi.OrderApproved:
		o.TicketID = evt.TicketID
		o.State = orderapi.Approved

	case *orderapi.OrderRejected:
		o.State = orderapi.Rejected

	case *OrderCancelling:
		o.State = orderapi.CancelPending

	case *OrderCancellingUndone:
		o.State = orderapi.Approved

	case *orderapi.OrderCancelled:
		o.State = orderapi.Cancelled

	case *orderapi.OrderProposedRevision:
		o.State = orderapi.RevisionPending

	case *OrderRevisingUndone:
		o.State = orderapi.Approved

	case *orderapi.OrderRevised:
		o.State = orderapi.Approved
		for menuItemID, quantity := range evt.RevisedQuantities {
			for idx, item := range o.LineItems {
				if item.MenuItemID == menuItemID {
					o.LineItems[idx].Quantity = quantity
				}
			}
		}

	default:
		return errs.NewError(fmt.Sprintf("unhandled event `%T`", event), ErrOrderUnhandledEvent)
	}

	return nil
}

func (o *Order) ApplySnapshot(snapshot core.Snapshot) error {
	switch ss := snapshot.(type) {
	case *OrderSnapshot:
		o.RestaurantID = ss.RestaurantID
		o.ConsumerID = ss.ConsumerID
		o.TicketID = ss.TicketID
		o.LineItems = ss.LineItems
		o.DeliverTo = ss.DeliverTo
		o.DeliverAt = ss.DeliverAt
		o.State = ss.State

	default:
		return errs.NewError(fmt.Sprintf("unhandled snapshot `%T`", snapshot), ErrOrderUnhandledSnapshot)
	}

	return nil
}

func (o *Order) ToSnapshot() (core.Snapshot, error) {
	return &OrderSnapshot{
		RestaurantID: o.RestaurantID,
		ConsumerID:   o.ConsumerID,
		TicketID:     o.TicketID,
		LineItems:    o.LineItems,
		DeliverTo:    o.DeliverTo,
		DeliverAt:    o.DeliverAt,
		State:        o.State,
	}, nil
}

func (o *Order) OrderTotal() int {
	total := 0
	for _, item := range o.LineItems {
		total += item.GetTotal()
	}

	return total
}

func (o *Order) RevisedOrderTotal(currentTotal int, revisedQuantities commonapi.MenuItemQuantities) int {
	delta := 0
	for menuItemID, quantity := range revisedQuantities {
		for _, lineItem := range o.LineItems {
			if lineItem.MenuItemID == menuItemID {
				delta += lineItem.Price * (quantity - lineItem.Quantity)
			}
		}
	}
	newTotal := currentTotal + delta

	return newTotal
}
