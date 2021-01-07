package queries

import (
	"time"
)

type OrderHistory struct {
	OrderID        string    `json:"order_id"`
	Status         string    `json:"status"`
	RestaurantID   string    `json:"restaurant_id"`
	RestaurantName string    `json:"restaurant_name"`
	CreatedAt      time.Time `json:"created_at"`
}
