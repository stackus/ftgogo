package domain

import (
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

var (
	ErrRestaurantNotFound = errors.Wrap(errors.ErrNotFound, "restaurant not found")
	ErrMenuItemNotFound   = errors.Wrap(errors.ErrNotFound, "menu item not found")
)

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
	return errors.ErrUnimplemented
}
