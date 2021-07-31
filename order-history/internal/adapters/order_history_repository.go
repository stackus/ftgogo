package adapters

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/order-history/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type OrderHistoryFilters struct {
	Since    time.Time           // rely on the .IsZero()
	Keywords []string            // ignored if empty
	Status   orderapi.OrderState // no pointer necessary; zero value == Unknown
	Next     string              // ignored if empty
	Limit    int                 // default to OrderHistoryLimit if not provided
}

// TODO update FindConsumerOrders to return a (*FindConsumerOrdersResult, error) pair

type OrderHistoryRepository interface {
	FindConsumerOrders(ctx context.Context, consumerID string, filters OrderHistoryFilters) ([]*domain.OrderHistory, string, error)
	Find(ctx context.Context, orderHistoryID string) (*domain.OrderHistory, error)
	Save(ctx context.Context, orderHistory *domain.OrderHistory) error
	UpdateStatus(ctx context.Context, orderHistoryID string, status orderapi.OrderState) error
	Update(ctx context.Context, orderHistoryID string, orderHistory *domain.OrderHistory) error
}
