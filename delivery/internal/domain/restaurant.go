package domain

import (
	"serviceapis/deliveryapi"
	"shared-go/errs"
)

// Restaurant errors
var (
	ErrRestaurantNotFound = errs.NewError("restaurant not found", errs.ErrNotFound)
)

type Restaurant struct {
	RestaurantID string
	Name         string
	Address      deliveryapi.Address
	// note: no menu items
}
