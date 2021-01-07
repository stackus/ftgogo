package domain

import (
	"github.com/stackus/edat/core"
	"serviceapis/commonapi"
	"serviceapis/restaurantapi"
	"shared-go/errs"
)

// Restaurant errors
var (
	ErrRestaurantNotFound = errs.NewError("restaurant not found", errs.ErrNotFound)
	ErrMenuItemNotFound   = errs.NewError("menu item not found", errs.ErrNotFound)
)

type Restaurant struct {
	core.EntityBase
	RestaurantID string
	Name         string
	Address      commonapi.Address
	MenuItems    []restaurantapi.MenuItem
}

var _ core.Entity = (*Restaurant)(nil)

func (Restaurant) EntityName() string {
	return "restaurantservice.Restaurant"
}

func (r Restaurant) ID() string {
	return r.RestaurantID
}

// FindMenuItem locates the local menu item record for a given menuItemID
func (r *Restaurant) FindMenuItem(menuItemID string) (restaurantapi.MenuItem, error) {
	for _, item := range r.MenuItems {
		if menuItemID == item.ID {
			return item, nil
		}
	}

	return restaurantapi.MenuItem{}, ErrMenuItemNotFound
}

// ReviseMenu updates the local menu
// NOT IMPLEMENTED
func (r *Restaurant) ReviseMenu([]restaurantapi.MenuItem) error {
	return errs.ErrNotImplemented
}
