package orderapi

import (
	"time"

	"github.com/stackus/ftgogo/serviceapis/commonapi"

	"github.com/stackus/edat/core"
)

func registerEvents() {
	core.RegisterEvents(
		OrderCreated{}, OrderApproved{}, OrderRejected{},
		OrderCancelled{},
		OrderProposedRevision{}, OrderRevised{},
	)
}

type OrderEvent struct{}

func (OrderEvent) DestinationChannel() string { return OrderAggregateChannel }

type OrderCreated struct {
	OrderEvent
	ConsumerID     string
	RestaurantID   string
	RestaurantName string
	LineItems      []LineItem
	OrderTotal     int
	DeliverAt      time.Time
	DeliverTo      commonapi.Address
}

func (OrderCreated) EventName() string { return "orderapi.OrderCreated" }

type OrderApproved struct {
	OrderEvent
	TicketID string
}

func (OrderApproved) EventName() string { return "orderapi.OrderApproved" }

type OrderRejected struct{ OrderEvent }

func (OrderRejected) EventName() string { return "orderapi.OrderRejected" }

type OrderProposedRevision struct {
	OrderEvent
	CurrentOrderTotal int
	NewOrderTotal     int
	Revisions         map[string]int
}

func (OrderProposedRevision) EventName() string { return "orderapi.OrderProposedRevision" }

type OrderCancelled struct{ OrderEvent }

func (OrderCancelled) EventName() string { return "orderapi.OrderCancelled" }

type OrderRevised struct {
	OrderEvent
	CurrentOrderTotal int
	NewOrderTotal     int
	RevisedQuantities map[string]int
}

func (OrderRevised) EventName() string { return "orderapi.OrderRevised" }
