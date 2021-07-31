package domain

import (
	"time"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
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
	PickUpAddress     *commonapi.Address
	DeliveryAddress   *commonapi.Address
	Status            DeliveryStatus
	PickUpTime        time.Time
	ReadyBy           time.Time
}

func (s DeliveryStatus) String() string {
	switch s {
	case DeliveryPending:
		return "PENDING"
	case DeliveryScheduled:
		return "SCHEDULED"
	case DeliveryCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}
