package domain

import (
	"github.com/google/uuid"
	"github.com/stackus/edat/core"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

// Restaurant errors
var (
	ErrRestaurantNotFound = errors.Wrap(errors.ErrNotFound, "restaurant not found")
	ErrMenuItemNotFound   = errors.Wrap(errors.ErrNotFound, "menu item not found")
)

type Restaurant struct {
	core.EntityBase
	RestaurantID string
	Name         string
	Address      *commonapi.Address
	MenuItems    []restaurantapi.MenuItem
}

var _ core.Entity = (*Restaurant)(nil)

func (Restaurant) EntityName() string {
	return "restaurantservice.Restaurant"
}

func (r Restaurant) ID() string {
	return r.RestaurantID
}

// CreateRestaurant builds a new Restaurant instance
func CreateRestaurant(name string, address *commonapi.Address, menuItems []restaurantapi.MenuItem) *Restaurant {
	r := &Restaurant{
		RestaurantID: uuid.New().String(),
		Name:         name,
		Address:      address,
		MenuItems:    menuItems,
	}

	r.AddEvent(&restaurantapi.RestaurantCreated{
		Name:    r.Name,
		Address: r.Address,
		Menu:    r.MenuItems,
	})

	return r
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
