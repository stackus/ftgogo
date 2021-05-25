package adapters

import (
	"context"

	"github.com/stackus/ftgogo/restaurant/internal/domain"
)

type restaurantPostgresPublisherMiddleware struct {
	domain.RestaurantRepository
	publisher domain.RestaurantPublisher
}

var _ domain.RestaurantRepository = (*restaurantPostgresPublisherMiddleware)(nil)

func NewRestaurantPostgresPublisherMiddleware(repository domain.RestaurantRepository, publisher domain.RestaurantPublisher) domain.RestaurantRepository {
	return &restaurantPostgresPublisherMiddleware{
		RestaurantRepository: repository,
		publisher:            publisher,
	}

}

func (r restaurantPostgresPublisherMiddleware) Save(ctx context.Context, restaurant *domain.Restaurant) error {
	err := r.RestaurantRepository.Save(ctx, restaurant)
	if err != nil {
		return err
	}

	return r.publisher.PublishEntityEvents(ctx, restaurant)
}

func (r restaurantPostgresPublisherMiddleware) Update(ctx context.Context, restaurantID string, restaurant *domain.Restaurant) error {
	err := r.RestaurantRepository.Update(ctx, restaurantID, restaurant)
	if err != nil {
		return err
	}

	return r.publisher.PublishEntityEvents(ctx, restaurant)
}