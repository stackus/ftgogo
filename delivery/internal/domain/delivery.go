package domain

import (
	"time"

	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/serviceapis/deliveryapi"
)

type DeliveryStatus int

const (
	DeliveryPending DeliveryStatus = iota
	DeliveryScheduled
	DeliveryCancelled
)

type Delivery struct {
	DeliveryID        string
	RestaurantID      string
	AssignedCourierID string
	PickUpAddress     deliveryapi.Address
	DeliveryAddress   deliveryapi.Address
	Status            DeliveryStatus
	PickUpTime        time.Time
	ReadyBy           time.Time
}

// Delivery errors
var (
	ErrDeliveryNotFound = errors.Wrap(errors.ErrNotFound, "delivery not found")
)

func (s DeliveryStatus) String() string {
	switch s {
	case DeliveryPending:
		return "PENDING"
	case DeliveryScheduled:
		return "SCHEDULED"
	case DeliveryCancelled:
		return "CANCELLED"
	}

	return "UNKNOWN"
}

func (d *Delivery) Schedule(readyBy time.Time, assignedCourierID string) {
	d.ReadyBy = readyBy
	d.AssignedCourierID = assignedCourierID
	d.Status = DeliveryScheduled
}

func (d *Delivery) Cancel() {
	d.AssignedCourierID = ""
	d.Status = DeliveryCancelled
}
