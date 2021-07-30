package queries

import (
	"context"

	"github.com/stackus/ftgogo/restaurant/internal/application/ports"
	"github.com/stackus/ftgogo/restaurant/internal/domain"
)

type GetRestaurant struct {
	RestaurantID string
}

type GetRestaurantHandler struct {
	repo ports.RestaurantRepository
}

func NewGetRestaurantHandler(repo ports.RestaurantRepository) GetRestaurantHandler {
	return GetRestaurantHandler{repo: repo}
}

func (h GetRestaurantHandler) Handle(ctx context.Context, query GetRestaurant) (*domain.Restaurant, error) {
	restaurant, err := h.repo.Find(ctx, query.RestaurantID)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}
