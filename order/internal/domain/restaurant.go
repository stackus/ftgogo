package domain

import (
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"

	"github.com/stackus/errors"
)

// Restaurant errors
var (
	ErrRestaurantNotFound = errors.Wrap(errors.ErrNotFound, "restaurant not found")
	ErrMenuItemNotFound   = errors.Wrap(errors.ErrNotFound, "menu item not found")
)

type Restaurant struct {
	RestaurantID string
	Name         string
	// note: no address
	MenuItems []restaurantapi.MenuItem
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
func (r *Restaurant) ReviseMenu(_ []restaurantapi.MenuItem) error {
	return errors.ErrUnimplemented
}
