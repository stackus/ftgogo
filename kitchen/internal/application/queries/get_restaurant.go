package queries

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type GetRestaurant struct {
	RestaurantID string
}

type GetRestaurantHandler struct {
	repo adapters.RestaurantRepository
}

func NewGetRestaurantHandler(restaurantRepo adapters.RestaurantRepository) GetRestaurantHandler {
	return GetRestaurantHandler{repo: restaurantRepo}
}

func (h GetRestaurantHandler) Handle(ctx context.Context, query GetRestaurant) (*domain.Restaurant, error) {
	restaurant, err := h.repo.Find(ctx, query.RestaurantID)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}
