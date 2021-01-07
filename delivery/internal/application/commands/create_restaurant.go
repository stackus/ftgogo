package commands

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/domain"
	"serviceapis/commonapi"
)

type CreateRestaurant struct {
	RestaurantID string
	Name         string
	Address      commonapi.Address
}

type CreateRestaurantHandler struct {
	repo domain.RestaurantRepository
}

func NewCreateRestaurantHandler(restaurantRepo domain.RestaurantRepository) CreateRestaurantHandler {
	return CreateRestaurantHandler{repo: restaurantRepo}
}

func (h CreateRestaurantHandler) Handle(ctx context.Context, cmd CreateRestaurant) error {
	return h.repo.Save(ctx, &domain.Restaurant{
		RestaurantID: cmd.RestaurantID,
		Name:         cmd.Name,
		Address:      cmd.Address,
	})
}
