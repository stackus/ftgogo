package domain

import (
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

type Restaurant struct {
	RestaurantID string
	Name         string
	Address      *commonapi.Address
	MenuItems    []restaurantapi.MenuItem
}
