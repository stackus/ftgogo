package queries

import (
	"context"

	"github.com/stackus/ftgogo/customer-web/internal/adapters"
	"github.com/stackus/ftgogo/customer-web/internal/domain"
)

type SearchOrders struct {
	ConsumerID string
	Filters    *domain.SearchOrdersFilters
	Next       string
	Limit      int
}

type SearchOrdersHandler struct {
	repo adapters.OrderHistoryRepository
}

func NewSearchOrdersHandler(repo adapters.OrderHistoryRepository) SearchOrdersHandler {
	return SearchOrdersHandler{repo: repo}
}

func (h SearchOrdersHandler) Handle(ctx context.Context, cmd SearchOrders) (*adapters.SearchOrdersResult, error) {
	return h.repo.SearchOrders(ctx, adapters.SearchOrders{
		ConsumerID: cmd.ConsumerID,
		Filters:    cmd.Filters,
		Next:       cmd.Next,
		Limit:      cmd.Limit,
	})
}
