package domain

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type Order struct {
	OrderID string
	Total   int
	Status  string
	// TODO EstimatedDelivery time.Duration
	// TODO other data
}

type (
	CreateOrder struct {
		ConsumerID   string
		RestaurantID string
		DeliverAt    time.Time
		DeliverTo    *commonapi.Address
		LineItems    commonapi.MenuItemQuantities
	}

	ReviseOrder struct {
		OrderID           string
		ConsumerID        string
		RevisedQuantities commonapi.MenuItemQuantities
	}
)

type OrderRepository interface {
	Create(ctx context.Context, createOrder CreateOrder) (string, error)
	Find(ctx context.Context, orderID string) (*Order, error)
	Revise(ctx context.Context, reviseOrder ReviseOrder) error
	Cancel(ctx context.Context, orderID string) error
}
