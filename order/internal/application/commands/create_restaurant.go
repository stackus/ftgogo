package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

type CreateRestaurant struct {
	RestaurantID string
	Name         string
	Menu         []restaurantapi.MenuItem
}

type CreateRestaurantHandler struct {
	repo ports.RestaurantRepository
}

func NewCreateRestaurantHandler(restaurantRepo ports.RestaurantRepository) CreateRestaurantHandler {
	return CreateRestaurantHandler{repo: restaurantRepo}
}

func (h CreateRestaurantHandler) Handle(ctx context.Context, cmd CreateRestaurant) error {
	return h.repo.Save(ctx, &domain.Restaurant{
		RestaurantID: cmd.RestaurantID,
		Name:         cmd.Name,
		MenuItems:    cmd.Menu,
	})
}
