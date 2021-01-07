package restaurantapi

import (
	"github.com/stackus/edat/core"
)

func registerCommands() {
	core.RegisterCommands(ValidateMenuItems{})
}

type RestaurantServiceCommand struct{}

func (RestaurantServiceCommand) DestinationChannel() string { return RestaurantServiceCommandChannel }

type ValidateMenuItems struct {
	RestaurantServiceCommand
	RestaurantID string
	MenuItems    []string
}

func (ValidateMenuItems) CommandName() string { return "restaurantapi.ValidateMenuItems" }
