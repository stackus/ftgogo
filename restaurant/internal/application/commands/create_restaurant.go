package commands

import (
	"context"

	"github.com/stackus/ftgogo/restaurant/internal/application/ports"
	"github.com/stackus/ftgogo/restaurant/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

type CreateRestaurant struct {
	Name      string
	Address   *commonapi.Address
	MenuItems []restaurantapi.MenuItem
}

type CreateRestaurantHandler struct {
	repo ports.RestaurantRepository
}

func NewCreateRestaurantHandler(restaurantRepo ports.RestaurantRepository) CreateRestaurantHandler {
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
