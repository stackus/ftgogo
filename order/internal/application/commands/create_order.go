package commands

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type CreateOrder struct {
	ConsumerID   string
	RestaurantID string
	DeliverAt    time.Time
	DeliverTo    *commonapi.Address
	LineItems    map[string]int
}

type CreateOrderHandler struct {
	orderRepo      ports.OrderRepository
	restaurantRepo ports.RestaurantRepository
}

func NewCreateOrderHandler(orderRepo ports.OrderRepository, restaurantRepo ports.RestaurantRepository) CreateOrderHandler {
	return CreateOrderHandler{
		orderRepo:      orderRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (h CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (string, error) {
	restaurant, err := h.restaurantRepo.Find(ctx, cmd.RestaurantID)
	if err != nil {
		return "", err
	}

	order, err := h.orderRepo.Save(ctx, &domain.CreateOrder{
		ConsumerID: cmd.ConsumerID,
		Restaurant: restaurant,
		LineItems:  cmd.LineItems,
		DeliverAt:  cmd.DeliverAt,
		DeliverTo:  cmd.DeliverTo,
	})
	if err != nil {
		return "", err
	}

	return order.ID(), nil
}
