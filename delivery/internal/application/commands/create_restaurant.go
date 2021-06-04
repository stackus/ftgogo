package commands

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/adapters"
	"github.com/stackus/ftgogo/delivery/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/deliveryapi"
)

type CreateRestaurant struct {
	RestaurantID string
	Name         string
	Address      deliveryapi.Address
}

type CreateRestaurantHandler struct {
	repo adapters.RestaurantRepository
}

func NewCreateRestaurantHandler(restaurantRepo adapters.RestaurantRepository) CreateRestaurantHandler {
	return CreateRestaurantHandler{repo: restaurantRepo}
}

func (h CreateRestaurantHandler) Handle(ctx context.Context, cmd CreateRestaurant) error {
	return h.repo.Save(ctx, &domain.Restaurant{
		RestaurantID: cmd.RestaurantID,
		Name:         cmd.Name,
		Address:      cmd.Address,
	})
}
