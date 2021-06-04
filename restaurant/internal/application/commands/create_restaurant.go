package commands

import (
	"context"

	"github.com/stackus/ftgogo/restaurant/internal/adapters"
	"github.com/stackus/ftgogo/restaurant/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

type CreateRestaurant struct {
	Name      string
	Address   restaurantapi.Address
	MenuItems []restaurantapi.MenuItem
}

type CreateRestaurantHandler struct {
	repo adapters.RestaurantRepository
}

func NewCreateRestaurantHandler(restaurantRepo adapters.RestaurantRepository) CreateRestaurantHandler {
	return CreateRestaurantHandler{
		repo: restaurantRepo,
	}
}

func (h CreateRestaurantHandler) Handle(ctx context.Context, cmd CreateRestaurant) (string, error) {
	restaurant := domain.CreateRestaurant(cmd.Name, cmd.Address, cmd.MenuItems)

	err := h.repo.Save(ctx, restaurant)
	if err != nil {
		return "", err
	}

	return restaurant.RestaurantID, nil
}
