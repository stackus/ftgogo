package adapters

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/customer-web/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
)

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
	Find(ctx context.Context, findOrder FindOrder) (*domain.Order, error)
	Revise(ctx context.Context, reviseOrder ReviseOrder) (orderapi.OrderState, error)
	Cancel(ctx context.Context, cancelOrder CancelOrder) (orderapi.OrderState, error)
}
