package domain

import (
	"time"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type OrderHistory struct {
	OrderID        string
	ConsumerID     string
	RestaurantID   string
	RestaurantName string
	LineItems      []orderapi.LineItem
	OrderTotal     int
	Status         orderapi.OrderState
	Keywords       []string
	CreatedAt      time.Time
}
