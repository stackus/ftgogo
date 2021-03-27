package domain

import (
	"time"

	"github.com/stackus/edat/core"

	"serviceapis/orderapi"
)

func registerOrderSnapshots() {
	core.RegisterSnapshots(OrderSnapshot{})
}

type OrderSnapshot struct {
	ConsumerID   string
	RestaurantID string
	TicketID     string
	LineItems    []orderapi.LineItem
	DeliverAt    time.Time
	DeliverTo    orderapi.Address
	State        orderapi.OrderState
}

func (OrderSnapshot) SnapshotName() string { return "orderservice.OrderSnapshot" }
