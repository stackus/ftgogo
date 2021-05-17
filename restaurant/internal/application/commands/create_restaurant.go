package commands

import (
	"context"

	"github.com/stackus/ftgogo/restaurant/internal/domain"
	"serviceapis/restaurantapi"
)

type CreateRestaurant struct {
	Name      string
	Address   restaurantapi.Address
	MenuItems []restaurantapi.MenuItem
}

type CreateRestaurantHandler struct {
	repo domain.RestaurantRepository
}

func NewCreateRestaurantHandler(restaurantRepo domain.RestaurantRepository) CreateRestaurantHandler {
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
