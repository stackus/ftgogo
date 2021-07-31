package domain

import (
	"time"

	"github.com/stackus/edat/core"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
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
	DeliverTo    *commonapi.Address
	State        orderapi.OrderState
}

func (OrderSnapshot) SnapshotName() string { return "orderservice.OrderSnapshot" }
