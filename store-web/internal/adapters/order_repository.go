package adapters

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/store-web/internal/domain"
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
	Find(ctx context.Context, findOrder FindOrder) (*domain.Order, error)
	Cancel(ctx context.Context, cancelOrder CancelOrder) (orderapi.OrderState, error)
}
