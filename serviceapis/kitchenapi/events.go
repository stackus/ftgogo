package kitchenapi

import (
	"time"

	"github.com/stackus/edat/core"
	"serviceapis/commonapi"
)

func registerEvents() {
	core.RegisterEvents(
		TicketCreated{}, TicketCancelled{}, TicketRevised{},
		TicketAccepted{},
	)
}

type TicketEvent struct{}

func (TicketEvent) DestinationChannel() string { return TicketAggregateChannel }

type TicketCreated struct {
	TicketEvent
	OrderID      string
	RestaurantID string
	LineItems    []LineItem
}

func (TicketCreated) EventName() string { return "kitchenapi.TicketCreated" }

type TicketCancelled struct {
	TicketEvent
	OrderID string
}

func (TicketCancelled) EventName() string { return "kitchenapi.TicketCancelled" }

type TicketRevised struct {
	TicketEvent
	RevisedQuantities commonapi.MenuItemQuantities
}

func (TicketRevised) EventName() string { return "kitchenapi.TicketRevised" }

type TicketAccepted struct {
	TicketEvent
	OrderID    string
	AcceptedAt time.Time
	ReadyBy    time.Time
}

func (TicketAccepted) EventName() string { return "kitchenapi.TicketAccepted" }
