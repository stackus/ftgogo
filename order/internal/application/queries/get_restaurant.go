package queries

import (
	"context"

	"github.com/stackus/ftgogo/order/internal/domain"
)

type GetRestaurant struct {
	RestaurantID string
}

type GetRestaurantHandler struct {
	repo domain.RestaurantRepository
}

func NewGetRestaurantHandler(restaurantRepo domain.RestaurantRepository) GetRestaurantHandler {
	return GetRestaurantHandler{repo: restaurantRepo}
}

func (h GetRestaurantHandler) Handle(ctx context.Context, query GetRestaurant) (*domain.Restaurant, error) {
	restaurant, err := h.repo.Find(ctx, query.RestaurantID)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}
