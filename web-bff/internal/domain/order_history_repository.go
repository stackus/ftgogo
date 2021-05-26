package domain

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type OrderHistory struct {
}

type SearchOrderHistoriesFilters struct {
	Keywords []string
	Since    time.Time
	Status   orderapi.OrderState
}

type (
	SearchOrderHistories struct {
		ConsumerID string
		Filters    SearchOrderHistoriesFilters
		Next       string
		Limit      int
	}
)

type OrderHistoryRepository interface {
	SearchOrderHistories(ctx context.Context, search SearchOrderHistories) ([]*OrderHistory, string, error)
}
