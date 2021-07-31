package commands

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/customer-web/internal/adapters"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type CreateOrder struct {
	ConsumerID   string
	RestaurantID string
	AddressID    string
	LineItems    commonapi.MenuItemQuantities
}

type CreateOrderHandler struct {
	orderRepo    adapters.OrderRepository
	consumerRepo adapters.ConsumerRepository
}

func NewCreateOrderHandler(orderRepo adapters.OrderRepository, consumerRepo adapters.ConsumerRepository) CreateOrderHandler {
	return CreateOrderHandler{
		orderRepo:    orderRepo,
		consumerRepo: consumerRepo,
	}
}

func (h CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (string, error) {
	address, err := h.consumerRepo.FindAddress(ctx, adapters.FindConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
	})
	if err != nil {
		return "", err
	}

	return h.orderRepo.Create(ctx, adapters.CreateOrder{
		ConsumerID:   cmd.ConsumerID,
		RestaurantID: cmd.RestaurantID,
		DeliverAt:    time.Now().Add(30 * time.Minute),
		DeliverTo:    address,
		LineItems:    cmd.LineItems,
	})
}
