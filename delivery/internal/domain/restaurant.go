package domain

import (
	"serviceapis/commonapi"
	"shared-go/errs"
)

// Restaurant errors
var (
	ErrRestaurantNotFound = errs.NewError("restaurant not found", errs.ErrNotFound)
)

type Restaurant struct {
	RestaurantID string
	Name         string
	Address      commonapi.Address
	// note: no menu items
}
