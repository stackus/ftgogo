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
	Status         orderapi.OrderState
	CreatedAt      time.Time
}

type SearchOrdersFilters struct {
	Keywords []string
	Since    time.Time
	Status   orderapi.OrderState
}
