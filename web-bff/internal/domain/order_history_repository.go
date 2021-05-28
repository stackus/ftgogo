package domain

import (
	"context"
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

type (
	SearchOrders struct {
		ConsumerID string
		Filters    *SearchOrdersFilters
		Next       string
		Limit      int
	}

	SearchOrdersResult struct {
		Orders []*OrderHistory
		Next   string
	}
)

type OrderHistoryRepository interface {
	SearchOrders(ctx context.Context, search SearchOrders) (*SearchOrdersResult, error)
}
