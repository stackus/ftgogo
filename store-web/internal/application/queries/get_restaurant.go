package queries

import (
	"context"

	"github.com/stackus/ftgogo/store-web/internal/adapters"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type GetRestaurant struct {
	RestaurantID string
}

type GetRestaurantHandler struct {
	repo adapters.RestaurantRepository
}

func NewGetRestaurantHandler(repo adapters.RestaurantRepository) GetRestaurantHandler {
	return GetRestaurantHandler{repo: repo}
}

func (h GetRestaurantHandler) Handle(ctx context.Context, query GetRestaurant) (*domain.Restaurant, error) {
	return h.repo.Find(ctx, adapters.FindRestaurant{RestaurantID: query.RestaurantID})
}
