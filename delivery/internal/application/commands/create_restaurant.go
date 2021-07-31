package commands

import (
	"context"

	"github.com/stackus/ftgogo/delivery/internal/application/ports"
	"github.com/stackus/ftgogo/delivery/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type CreateRestaurant struct {
	RestaurantID string
	Name         string
	Address      *commonapi.Address
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
		Address:      cmd.Address,
	})
}
