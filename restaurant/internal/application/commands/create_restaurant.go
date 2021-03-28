package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/stackus/ftgogo/restaurant/internal/domain"
	"serviceapis/restaurantapi"
)

type CreateRestaurant struct {
	Name      string
	Address   restaurantapi.Address
	MenuItems []restaurantapi.MenuItem
}

type CreateRestaurantHandler struct {
	repo      domain.RestaurantRepository
	publisher domain.RestaurantPublisher
}

func NewCreateRestaurantHandler(restaurantRepo domain.RestaurantRepository, restaurantPublisher domain.RestaurantPublisher) CreateRestaurantHandler {
	return CreateRestaurantHandler{
		repo:      restaurantRepo,
		publisher: restaurantPublisher,
	}
}

func (h CreateRestaurantHandler) Handle(ctx context.Context, cmd CreateRestaurant) (string, error) {
	restaurant := &domain.Restaurant{
		RestaurantID: uuid.New().String(),
		Name:         cmd.Name,
		Address:      cmd.Address,
		MenuItems:    cmd.MenuItems,
	}

	restaurant.AddEvent(&restaurantapi.RestaurantCreated{
		Name:    cmd.Name,
		Address: cmd.Address,
		Menu:    cmd.MenuItems,
	})

	err := h.repo.Save(ctx, restaurant)
	if err != nil {
		return "", err
	}

	err = h.publisher.PublishEntityEvents(ctx, restaurant)
	if err != nil {
		return "", err
	}

	return restaurant.RestaurantID, nil
}
