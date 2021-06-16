package restaurantapi

import (
	"github.com/stackus/edat/core"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

func registerEvents() {
	core.RegisterEvents(RestaurantCreated{}, RestaurantMenuRevised{})
}

type RestaurantEvent struct{}

func (RestaurantEvent) DestinationChannel() string { return RestaurantAggregateChannel }

type RestaurantCreated struct {
	RestaurantEvent
	Name    string
	Address *commonapi.Address
	Menu    []MenuItem
}

func (RestaurantCreated) EventName() string { return "restaurantapi.RestaurantCreated" }

type RestaurantMenuRevised struct {
	RestaurantEvent
	Menu []MenuItem
}

func (RestaurantMenuRevised) EventName() string { return "restaurantapi.RestaurantMenuRevised" }
