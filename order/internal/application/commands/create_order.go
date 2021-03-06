package commands

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/stackus/ftgogo/order/internal/domain"
	"serviceapis/commonapi"
	"serviceapis/orderapi"
)

type CreateOrder struct {
	ConsumerID   string
	RestaurantID string
	DeliverAt    time.Time
	DeliverTo    commonapi.Address
	LineItems    commonapi.MenuItemQuantities
}

type CreateOrderHandler struct {
	orderRepo      domain.OrderRepository
	restaurantRepo domain.RestaurantRepository
	publisher      domain.OrderPublisher
	logger         zerolog.Logger // TODO Eliminate passing logger everywhere use a pkg logger i.e. logging.Error()...
}

func NewCreateOrderHandler(
	orderRepo domain.OrderRepository,
	restaurantRepo domain.RestaurantRepository,
	orderPublisher domain.OrderPublisher,
	logger zerolog.Logger,
) CreateOrderHandler {
	return CreateOrderHandler{
		orderRepo:      orderRepo,
		restaurantRepo: restaurantRepo,
		publisher:      orderPublisher,
		logger:         logger,
	}
}

func (h CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (string, error) {
	restaurant, err := h.restaurantRepo.Find(ctx, cmd.RestaurantID)
	if err != nil {
		return "", err
	}

	total := 0
	orderLineItems := make([]orderapi.LineItem, 0, len(cmd.LineItems))
	for menuItemID, quantity := range cmd.LineItems {
		menuItem, mErr := restaurant.FindMenuItem(menuItemID)
		if mErr != nil {
			return "", mErr
		}
		item := orderapi.LineItem{
			MenuItemID: menuItemID,
			Name:       menuItem.Name,
			Price:      menuItem.Price,
			Quantity:   quantity,
		}
		total += item.GetTotal()
		orderLineItems = append(orderLineItems, item)
	}

	root, err := h.orderRepo.Save(ctx, &domain.CreateOrder{
		ConsumerID:     cmd.ConsumerID,
		RestaurantID:   cmd.RestaurantID,
		RestaurantName: restaurant.Name,
		LineItems:      orderLineItems,
		OrderTotal:     total,
		DeliverAt:      cmd.DeliverAt,
		DeliverTo:      cmd.DeliverTo,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("error while saving the order")
		return "", err
	}

	err = h.publisher.PublishEntityEvents(ctx, root)
	if err != nil {
		h.logger.Error().Err(err).Msg("error while publishing the message")
	}

	return root.AggregateID(), nil
}
