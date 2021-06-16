package commands

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
	"github.com/stackus/ftgogo/store-web/internal/adapters"
)

type CreateRestaurant struct {
	Name      string
	Address   *commonapi.Address
	MenuItems []restaurantapi.MenuItem
}

type CreateRestaurantHandler struct {
	repo adapters.RestaurantRepository
}

func NewCreateRestaurantHandler(repo adapters.RestaurantRepository) CreateRestaurantHandler {
	return CreateRestaurantHandler{repo: repo}
}

func (h CreateRestaurantHandler) Handle(ctx context.Context, cmd CreateRestaurant) (string, error) {
	return h.repo.Create(ctx, adapters.CreateRestaurant{
		Name:      cmd.Name,
		Address:   cmd.Address,
		MenuItems: cmd.MenuItems,
	})
}
