package domain

import (
	"serviceapis/restaurantapi"
	"shared-go/errs"
)

// Restaurant errors
var (
	ErrRestaurantNotFound = errs.NewError("restaurant not found", errs.ErrNotFound)
	ErrMenuItemNotFound   = errs.NewError("menu item not found", errs.ErrNotFound)
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
	return errs.ErrNotImplemented
}
