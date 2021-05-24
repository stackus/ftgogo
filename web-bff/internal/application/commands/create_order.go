package commands

import (
	"context"
	"time"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type CreateOrder struct {
	ConsumerID   string
	RestaurantID string
	AddressID    string
	LineItems    commonapi.MenuItemQuantities
}

type CreateOrderHandler struct {
	orderRepo    domain.OrderRepository
	consumerRepo domain.ConsumerRepository
}

func NewCreateOrderHandler(orderRepo domain.OrderRepository, consumerRepo domain.ConsumerRepository) CreateOrderHandler {
	return CreateOrderHandler{
		orderRepo:    orderRepo,
		consumerRepo: consumerRepo,
	}
}

func (h CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (string, error) {
	address, err := h.consumerRepo.FindAddress(ctx, domain.FindConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
	})
	if err != nil {
		return "", err
	}

	return h.orderRepo.Create(ctx, domain.CreateOrder{
		ConsumerID:   cmd.ConsumerID,
		RestaurantID: cmd.RestaurantID,
		DeliverAt:    time.Now().Add(30 * time.Minute),
		DeliverTo:    address,
		LineItems:    cmd.LineItems,
	})
}
