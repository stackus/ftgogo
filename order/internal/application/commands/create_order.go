package commands

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type CreateOrder struct {
	ConsumerID   string
	RestaurantID string
	DeliverAt    time.Time
	DeliverTo    commonapi.Address
	LineItems    map[string]int
}

type CreateOrderHandler struct {
	orderRepo      adapters.OrderRepository
	restaurantRepo adapters.RestaurantRepository
	logger         zerolog.Logger // TODO Eliminate passing logger everywhere use a pkg logger i.e. logging.Error()...
}

func NewCreateOrderHandler(orderRepo adapters.OrderRepository, restaurantRepo adapters.RestaurantRepository, logger zerolog.Logger) CreateOrderHandler {
	return CreateOrderHandler{
		orderRepo:      orderRepo,
		restaurantRepo: restaurantRepo,
		logger:         logger,
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
		h.logger.Error().Err(err).Msg("error while saving the order")
		return "", err
	}

	return order.ID(), nil
}
