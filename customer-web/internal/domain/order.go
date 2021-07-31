package domain

import (
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type Order struct {
	OrderID        string
	ConsumerID     string
	RestaurantID   string
	RestaurantName string
	Total          int
	Status         orderapi.OrderState
	// TODO EstimatedDelivery time.Duration
	// TODO other data
}
