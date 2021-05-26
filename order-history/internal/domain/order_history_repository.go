package domain

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

const OrderHistoryLimit = 20
const OrderHistoryMinimum = 1
const OrderHistoryMaximum = 50

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

type OrderHistoryFilters struct {
	Since    time.Time           // rely on the .IsZero()
	Keywords []string            // ignored if empty
	Status   orderapi.OrderState // no pointer necessary; zero value == Unknown
	Next     string              // ignored if empty
	Limit    int                 // default to OrderHistoryLimit if not provided
}

type OrderHistoryRepository interface {
	FindConsumerOrders(ctx context.Context, consumerID string, filters OrderHistoryFilters) ([]*OrderHistory, string, error)
	Find(ctx context.Context, orderHistoryID string) (*OrderHistory, error)
	Save(ctx context.Context, orderHistory *OrderHistory) error
	UpdateStatus(ctx context.Context, orderHistoryID string, status orderapi.OrderState) error
	Update(ctx context.Context, orderHistoryID string, orderHistory *OrderHistory) error
}
