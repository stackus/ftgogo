package queries

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/order-history/internal/adapters"
	"github.com/stackus/ftgogo/order-history/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type SearchOrderHistories struct {
	ConsumerID string
	Filter     *OrderHistoryFilters
	Next       string
	Limit      int
}

type SearchOrderHistoriesHandler struct {
	repo adapters.OrderHistoryRepository
}

type OrderHistoryFilters struct {
	Since    time.Time
	Keywords []string
	Status   orderapi.OrderState
}

type SearchOrderHistoriesResponse struct {
	Orders []*domain.OrderHistory
	Next   string
}

func NewSearchOrderHistoriesHandler(orderHistoryRepo adapters.OrderHistoryRepository) SearchOrderHistoriesHandler {
	return SearchOrderHistoriesHandler{repo: orderHistoryRepo}
}

func (h SearchOrderHistoriesHandler) Handle(ctx context.Context, query SearchOrderHistories) (*SearchOrderHistoriesResponse, error) {
	filters := adapters.OrderHistoryFilters{}

	if query.Next != "" {
		filters.Next = query.Next
	}

	filters.Limit = domain.OrderHistoryLimit
	if query.Limit >= domain.OrderHistoryMinimum && query.Limit <= domain.OrderHistoryMaximum {
		filters.Limit = query.Limit
	}

	if query.Filter != nil {
		filters.Keywords = query.Filter.Keywords

		if query.Filter.Status != orderapi.UnknownOrderState {
			filters.Status = query.Filter.Status
		}

		if !query.Filter.Since.IsZero() {
			filters.Since = query.Filter.Since
		}
	}

	orders, next, err := h.repo.FindConsumerOrders(ctx, query.ConsumerID, filters)
	if err != nil {
		return nil, err
	}

	return &SearchOrderHistoriesResponse{
		Orders: orders,
		Next:   next,
	}, nil
}
