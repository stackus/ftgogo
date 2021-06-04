package adapters

import (
	"context"

	"github.com/stackus/ftgogo/customer-web/internal/domain"
)

type (
	SearchOrders struct {
		ConsumerID string
		Filters    *domain.SearchOrdersFilters
		Next       string
		Limit      int
	}

	SearchOrdersResult struct {
		Orders []*domain.OrderHistory
		Next   string
	}
)

type OrderHistoryRepository interface {
	SearchOrders(ctx context.Context, search SearchOrders) (*SearchOrdersResult, error)
}
