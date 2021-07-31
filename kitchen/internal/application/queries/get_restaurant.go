package queries

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type GetRestaurant struct {
	RestaurantID string
}

type GetRestaurantHandler struct {
	repo ports.RestaurantRepository
}

func NewGetRestaurantHandler(restaurantRepo ports.RestaurantRepository) GetRestaurantHandler {
	return GetRestaurantHandler{repo: restaurantRepo}
}

func (h GetRestaurantHandler) Handle(ctx context.Context, query GetRestaurant) (*domain.Restaurant, error) {
	restaurant, err := h.repo.Find(ctx, query.RestaurantID)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}
