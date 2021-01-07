package domain

import (
	"github.com/stackus/edat/core"
	"serviceapis/kitchenapi"
)

func registerTicketEvents() {
	core.RegisterEvents(
		TicketCreated{}, TicketCreateConfirmed{}, TicketCreateCancelled{},
		TicketCancelling{}, TicketCancelUndone{}, TicketCancelConfirmed{},
		TicketRevising{}, TicketReviseConfirmed{}, TicketReviseUndone{},
	)
}

type TicketCreated struct {
	OrderID      string
	RestaurantID string
	LineItems    []kitchenapi.LineItem
}

func (TicketCreated) EventName() string { return "kitchenservice.TicketCreated" }

type TicketCreateConfirmed struct{}

func (TicketCreateConfirmed) EventName() string { return "kitchenservice.TicketCreateConfirmed" }

type TicketCreateCancelled struct{}

func (TicketCreateCancelled) EventName() string { return "kitchenservice.TicketCreateCancelled" }

type TicketCancelling struct{}

func (TicketCancelling) EventName() string { return "ktichenservice.TicketCancelling" }

type TicketCancelConfirmed struct{}

func (TicketCancelConfirmed) EventName() string { return "kitchenservice.TicketCancelConfirmed" }

type TicketCancelUndone struct{}

func (TicketCancelUndone) EventName() string { return "kitchenservice.TicketCancelUndone" }

type TicketRevising struct{}

func (TicketRevising) EventName() string { return "ktichenservice.TicketRevising" }

type TicketReviseConfirmed struct{}

func (TicketReviseConfirmed) EventName() string { return "kitchenservice.TicketReviseConfirmed" }

type TicketReviseUndone struct{}

func (TicketReviseUndone) EventName() string { return "kitchenservice.TicketReviseUndone" }
