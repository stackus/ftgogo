package domain

import (
	"fmt"
	"time"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"serviceapis/kitchenapi"
	"shared-go/errs"
)

var ErrTicketUnhandledCommand = errs.NewError("unhandled command in ticket aggregate", errs.ErrServerError)
var ErrTicketUnhandledEvent = errs.NewError("unhandled event in ticket aggregate", errs.ErrServerError)
var ErrTicketUnhandledSnapshot = errs.NewError("unhandled snapshot in ticket aggregate", errs.ErrServerError)

var ErrTicketInvalidState = errs.NewError("ticket state does not allow action", errs.ErrConflict)
var ErrTicketReadyByBeforeNow = errs.NewError("readyBy is not after now", errs.ErrNotAcceptable)

type TicketState int

const (
	CreatePending TicketState = iota
	AwaitingAcceptance
	Accepted
	Preparing
	ReadyForPickup
	PickedUp
	CancelPending
	Cancelled
	RevisionPending
)

type Ticket struct {
	es.AggregateBase
	OrderID          string
	RestaurantID     string
	LineItems        []kitchenapi.LineItem
	ReadyBy          time.Time
	AcceptedAt       time.Time
	PreparingTime    time.Time
	ReadyForPickUpAt time.Time
	PickedUpAt       time.Time
	State            TicketState
	PreviousState    TicketState
}

func NewTicket() es.Aggregate {
	return &Ticket{}
}

func (Ticket) EntityName() string {
	return "kitchenservice.Ticket"
}

func (t *Ticket) ProcessCommand(command core.Command) error {
	switch cmd := command.(type) {
	case *CreateTicket:
		t.AddEvent(&TicketCreated{
			OrderID:      cmd.OrderID,
			RestaurantID: cmd.RestaurantID,
			LineItems:    cmd.LineItems,
		})

	case *ConfirmCreateTicket:
		if t.State != CreatePending {
			return ErrTicketInvalidState
		}
		t.AddEvent(&kitchenapi.TicketCreated{
			OrderID:      t.OrderID,
			RestaurantID: t.RestaurantID,
			LineItems:    t.LineItems,
		})

	case *CancelCreateTicket:
		// Originally not implemented in ftgo-kitchen-service
		if t.State != CreatePending {
			return ErrTicketInvalidState
		}
		t.AddEvent(&TicketCreateCancelled{})

	case *CancelTicket:
		if t.State != AwaitingAcceptance && t.State != Accepted {
			return ErrTicketInvalidState
		}
		t.AddEvent(&TicketCancelling{})

	case *ConfirmCancelTicket:
		if t.State != CancelPending {
			return ErrTicketInvalidState
		}
		t.AddEvent(&kitchenapi.TicketCancelled{
			OrderID: t.OrderID,
		})

	case *UndoCancelTicket:
		if t.State != CancelPending {
			return ErrTicketInvalidState
		}
		t.AddEvent(&TicketCancelUndone{})

	case *ReviseTicket:
		if t.State != AwaitingAcceptance && t.State != Accepted {
			return ErrTicketInvalidState
		}
		t.AddEvent(&TicketRevising{})

	case *ConfirmReviseTicket:
		if t.State != RevisionPending {
			return ErrTicketInvalidState
		}
		t.AddEvent(&kitchenapi.TicketRevised{
			RevisedQuantities: cmd.RevisedQuantities,
		})

	case *UndoReviseTicket:
		if t.State != RevisionPending {
			return ErrTicketInvalidState
		}
		t.AddEvent(&TicketReviseUndone{})

	case *AcceptTicket:
		if t.State != AwaitingAcceptance {
			return ErrTicketInvalidState
		}
		if time.Now().After(cmd.ReadyBy) {
			return ErrTicketReadyByBeforeNow
		}

		t.AddEvent(&kitchenapi.TicketAccepted{
			OrderID:    t.OrderID,
			AcceptedAt: time.Now(),
			ReadyBy:    cmd.ReadyBy,
		})

	default:
		return errs.NewError(fmt.Sprintf("unhandled command `%T`", command), ErrTicketUnhandledCommand)
	}

	return nil
}

func (t *Ticket) ApplyEvent(event core.Event) error {
	switch evt := event.(type) {
	case *TicketCreated:
		t.OrderID = evt.OrderID
		t.RestaurantID = evt.RestaurantID
		t.LineItems = evt.LineItems
		t.State = CreatePending

	case *kitchenapi.TicketCreated:
		t.State = AwaitingAcceptance

	case *TicketCreateCancelled:
		t.State = Cancelled // possibly; not implemented in ftgo-kitchen-service

	case *TicketCancelling:
		t.PreviousState = t.State
		t.State = CancelPending

	case *kitchenapi.TicketCancelled:
		t.State = Cancelled

	case *TicketCancelUndone:
		t.State = t.PreviousState

	case *TicketRevising:
		t.PreviousState = t.State
		t.State = RevisionPending

	case *kitchenapi.TicketRevised:
		t.State = t.PreviousState

	case *TicketReviseUndone:
		t.State = t.PreviousState

	case *kitchenapi.TicketAccepted:
		t.AcceptedAt = evt.AcceptedAt
		t.ReadyBy = evt.ReadyBy
		t.State = Accepted // assume that this is the case; doesn't appear to be ever set in ftgo-kitchen-service

	default:
		return errs.NewError(fmt.Sprintf("unhandled event `%T`", event), ErrTicketUnhandledEvent)
	}

	return nil
}

func (t *Ticket) ApplySnapshot(snapshot core.Snapshot) error {
	switch ss := snapshot.(type) {
	case *TicketSnapshot:
		t.OrderID = ss.OrderID
		t.RestaurantID = ss.RestaurantID
		t.LineItems = ss.LineItems
		t.ReadyBy = ss.ReadyBy
		t.AcceptedAt = ss.AcceptedAt
		t.PreparingTime = ss.PreparingTime
		t.ReadyForPickUpAt = ss.ReadyForPickUpAt
		t.PickedUpAt = ss.PickedUpAt
		t.State = ss.State

	default:
		return errs.NewError(fmt.Sprintf("unhandled snapshot `%T`", snapshot), ErrTicketUnhandledSnapshot)
	}

	return nil
}

func (t *Ticket) ToSnapshot() (core.Snapshot, error) {
	return &TicketSnapshot{
		OrderID:          t.OrderID,
		RestaurantID:     t.RestaurantID,
		LineItems:        t.LineItems,
		ReadyBy:          t.ReadyBy,
		AcceptedAt:       t.AcceptedAt,
		PreparingTime:    t.PreparingTime,
		ReadyForPickUpAt: t.ReadyForPickUpAt,
		PickedUpAt:       t.PickedUpAt,
		State:            t.State,
	}, nil
}
