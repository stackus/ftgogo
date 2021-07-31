package adapters

import (
	"context"

	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

type RestaurantInmemRepository struct {
	restaurants map[string]*domain.Restaurant
}

var _ ports.RestaurantRepository = (*RestaurantInmemRepository)(nil)

func NewRestaurantInmemRepository() *RestaurantInmemRepository {
	return &RestaurantInmemRepository{restaurants: map[string]*domain.Restaurant{}}
}

func (r *RestaurantInmemRepository) Find(ctx context.Context, restaurantID string) (*domain.Restaurant, error) {
	if restaurant, exists := r.restaurants[restaurantID]; !exists {
		return nil, domain.ErrRestaurantNotFound
	} else {
		return restaurant, nil
	}
}

func (r *RestaurantInmemRepository) Save(ctx context.Context, restaurant *domain.Restaurant) error {
	if _, exists := r.restaurants[restaurant.RestaurantID]; exists {
		return errors.Wrap(errors.ErrConflict, "restaurant already exists")
	}
	r.restaurants[restaurant.RestaurantID] = restaurant
	return nil
}

func (r *RestaurantInmemRepository) Update(ctx context.Context, restaurantID string, restaurant *domain.Restaurant) error {
	if _, exists := r.restaurants[restaurantID]; !exists {
		return domain.ErrRestaurantNotFound
	}
	r.restaurants[restaurantID] = restaurant
	return nil
}
