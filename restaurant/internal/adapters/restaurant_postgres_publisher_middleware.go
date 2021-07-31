package adapters

import (
	"context"

	"github.com/stackus/ftgogo/restaurant/internal/application/ports"
	"github.com/stackus/ftgogo/restaurant/internal/domain"
)

type restaurantPostgresPublisherMiddleware struct {
	ports.RestaurantRepository
	publisher ports.RestaurantPublisher
}

var _ ports.RestaurantRepository = (*restaurantPostgresPublisherMiddleware)(nil)

func NewRestaurantPostgresPublisherMiddleware(repository ports.RestaurantRepository, publisher ports.RestaurantPublisher) ports.RestaurantRepository {
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
