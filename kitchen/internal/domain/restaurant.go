package domain

import (
	"serviceapis/restaurantapi"
	"shared-go/errs"
)

var ErrRestaurantNotFound = errs.NewError("restaurant not found", errs.ErrNotFound)
var ErrMenuItemNotFound = errs.NewError("menu item not found", errs.ErrNotFound)

type Restaurant struct {
	RestaurantID string
	Name         string
	MenuItems    []restaurantapi.MenuItem
}

func (r *Restaurant) FindMenuItem(menuItemID string) (restaurantapi.MenuItem, error) {
	for _, item := range r.MenuItems {
		if menuItemID == item.ID {
			return item, nil
		}
	}

	return restaurantapi.MenuItem{}, ErrMenuItemNotFound
}

func (r *Restaurant) ReviseMenu([]restaurantapi.MenuItem) error {
	return errs.ErrNotImplemented
}
