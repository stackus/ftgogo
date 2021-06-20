package ports

import (
	"context"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type RestaurantRepository interface {
	Find(ctx context.Context, restaurantID string) (*domain.Restaurant, error)
	Save(ctx context.Context, restaurant *domain.Restaurant) error
	Update(ctx context.Context, restaurantID string, restaurant *domain.Restaurant) error
}
