package steps

import (
	"github.com/google/uuid"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

type consumer struct {
	ID      string
	Address *commonapi.Address
}

var consumers = map[string]*consumer{
	"Able Anders": {
		ID: uuid.New().String(),
		Address: &commonapi.Address{
			Street1: "123 Main St.",
			Street2: "",
			City:    "SomeTown",
			State:   "ST",
			Zip:     "90210",
		},
	},
}

func getConsumerFromFixture(consumerName string) (*consumer, error) {
	consumer, exists := consumers[consumerName]
	if !exists {
		return nil, errors.Wrapf(errors.ErrNotFound, "no consumer '%s' exists in fixture data", consumerName)
	}

	return consumer, nil
}

var restaurants = map[string]*domain.Restaurant{
	"Best Foods": {
		RestaurantID: uuid.New().String(),
		Name:         "Best Foods",
		MenuItems: []restaurantapi.MenuItem{
			{
				ID:    "yummy",
				Name:  "Yummy Dish",
				Price: 999,
			},
			{
				ID:    "soup",
				Name:  "Hot Soup",
				Price: 599,
			},
			{
				ID:    "salad",
				Name:  "Chef Salad",
				Price: 599,
			},
		},
	},
	"Other Foods": {
		RestaurantID: uuid.New().String(),
		Name:         "Other Foods",
		MenuItems: []restaurantapi.MenuItem{
			{
				ID:    "yummy",
				Name:  "Yummy Dish",
				Price: 999,
			},
			{
				ID:    "soup",
				Name:  "Hot Soup",
				Price: 599,
			},
			{
				ID:    "salad",
				Name:  "Chef Salad",
				Price: 599,
			},
		},
	},
}

func getRestaurantFromFixture(restaurantName string) (*domain.Restaurant, error) {
	restaurant, exists := restaurants[restaurantName]
	if !exists {
		return nil, errors.Wrapf(errors.ErrNotFound, "no restaurant '%s' exists in fixture data", restaurantName)
	}

	return restaurant, nil
}
