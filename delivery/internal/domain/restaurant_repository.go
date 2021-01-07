package domain

import (
	"context"
)

type RestaurantRepository interface {
	Find(ctx context.Context, restaurantID string) (*Restaurant, error)
	Save(ctx context.Context, restaurant *Restaurant) error
	Update(ctx context.Context, restaurantID string, restaurant *Restaurant) error
}
