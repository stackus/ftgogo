package adapters

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type (
	CreateRestaurant struct {
		Name      string
		Address   *commonapi.Address
		MenuItems []restaurantapi.MenuItem
	}

	FindRestaurant struct {
		RestaurantID string
	}
)

type RestaurantRepository interface {
	Create(context.Context, CreateRestaurant) (string, error)
	Find(ctx context.Context, restaurant FindRestaurant) (*domain.Restaurant, error)
}
