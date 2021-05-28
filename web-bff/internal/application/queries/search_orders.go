package queries

import (
	"context"

	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type SearchOrders struct {
	ConsumerID string
	Filters    *domain.SearchOrdersFilters
	Next       string
	Limit      int
}

type SearchOrdersHandler struct {
	repo domain.OrderHistoryRepository
}

func NewSearchOrdersHandler(repo domain.OrderHistoryRepository) SearchOrdersHandler {
	return SearchOrdersHandler{repo: repo}
}

func (h SearchOrdersHandler) Handle(ctx context.Context, cmd SearchOrders) (*domain.SearchOrdersResult, error) {
	return h.repo.SearchOrders(ctx, domain.SearchOrders{
		ConsumerID: cmd.ConsumerID,
		Filters:    cmd.Filters,
		Next:       cmd.Next,
		Limit:      cmd.Limit,
	})
}
