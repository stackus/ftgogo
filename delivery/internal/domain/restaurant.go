package domain

import (
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

// Restaurant errors
var (
	ErrRestaurantNotFound = errors.Wrap(errors.ErrNotFound, "restaurant not found")
)

type Restaurant struct {
	RestaurantID string
	Name         string
	Address      *commonapi.Address
	// note: no menu items
}
