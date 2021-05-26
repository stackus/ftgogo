package domain

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

type Order struct {
	OrderID        string
	ConsumerID     string
	RestaurantID   string
	RestaurantName string
	Total          int
	Status         orderapi.OrderState
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

	FindOrder struct {
		OrderID string
	}

	CancelOrder FindOrder

	ReviseOrder struct {
		OrderID           string
		RevisedQuantities commonapi.MenuItemQuantities
	}
)

type OrderRepository interface {
	Create(ctx context.Context, createOrder CreateOrder) (string, error)
	Find(ctx context.Context, findOrder FindOrder) (*Order, error)
	Revise(ctx context.Context, reviseOrder ReviseOrder) (orderapi.OrderState, error)
	Cancel(ctx context.Context, cancelOrder CancelOrder) (orderapi.OrderState, error)
}
