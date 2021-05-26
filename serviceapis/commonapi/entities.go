package commonapi

import (
	commonpb "github.com/stackus/ftgogo/serviceapis/commonapi/pb"
)

// Generate spec imports for openapi documentation
//go:generate oapi-codegen -generate spec -o spec.gen.go -package commonapi  .\spec.yaml

type Address struct {
	Street1 string
	Street2 string
	City    string
	State   string
	Zip     string
}

func ToAddressProto(address *Address) *commonpb.Address {
	return &commonpb.Address{
		Street1: address.Street1,
		Street2: address.Street2,
		City:    address.City,
		State:   address.State,
		Zip:     address.Zip,
	}
}

func FromAddressProto(address *commonpb.Address) *Address {
	return &Address{
		Street1: address.Street1,
		Street2: address.Street2,
		City:    address.City,
		State:   address.State,
		Zip:     address.Zip,
	}
}

type MenuItemQuantities map[string]int

func ToMenuItemQuantitiesProto(quantities MenuItemQuantities) *commonpb.MenuItemQuantities {
	lineItems := make(map[string]int64, len(quantities))
	for itemID, qty := range quantities {
		lineItems[itemID] = int64(qty)
	}

	return &commonpb.MenuItemQuantities{Items: lineItems}
}

func FromMenuItemQuantitiesProto(quantities *commonpb.MenuItemQuantities) MenuItemQuantities {
	lineItems := make(MenuItemQuantities, len(quantities.Items))
	for itemID, qty := range quantities.Items {
		lineItems[itemID] = int(qty)
	}

	return lineItems
}
