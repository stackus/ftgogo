package domain

import (
	"time"

	"github.com/stackus/edat/core"

	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
)

func registerTicketSnapshots() {
	core.RegisterSnapshots(TicketSnapshot{})
}

type TicketSnapshot struct {
	OrderID          string
	RestaurantID     string
	LineItems        []kitchenapi.LineItem
	ReadyBy          time.Time
	AcceptedAt       time.Time
	PreparingTime    time.Time
	ReadyForPickUpAt time.Time
	PickedUpAt       time.Time
	State            TicketState
}

func (TicketSnapshot) SnapshotName() string { return "kitchenservice.TicketSnapshot" }
