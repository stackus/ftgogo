package domain

import (
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/serviceapis/deliveryapi"
)

// Restaurant errors
var (
	ErrRestaurantNotFound = errors.Wrap(errors.ErrNotFound, "restaurant not found")
)

type Restaurant struct {
	RestaurantID string
	Name         string
	Address      deliveryapi.Address
	// note: no menu items
}
