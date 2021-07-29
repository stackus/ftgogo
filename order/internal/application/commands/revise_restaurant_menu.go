package commands

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/application/ports"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

type ReviseRestaurantMenu struct {
	RestaurantID string
	Menu         []restaurantapi.MenuItem
}

type ReviseRestaurantMenuHandler struct {
	repo ports.RestaurantRepository
}

func NewReviseRestaurantMenuHandler(restaurantRepo ports.RestaurantRepository) ReviseRestaurantMenuHandler {
	return ReviseRestaurantMenuHandler{repo: restaurantRepo}
}

func (h ReviseRestaurantMenuHandler) Handle(ctx context.Context, cmd ReviseRestaurantMenu) error {
	restaurant, err := h.repo.Find(ctx, cmd.RestaurantID)
	if err != nil {
		return err
	}

	err = restaurant.ReviseMenu(cmd.Menu)
	if err != nil {
		return err
	}

	err = h.repo.Update(ctx, cmd.RestaurantID, restaurant)
	if err != nil {
		return err
	}

	return nil
}
