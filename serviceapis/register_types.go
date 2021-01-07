package serviceapis

import (
	"serviceapis/accountingapi"
	"serviceapis/consumerapi"
	"serviceapis/kitchenapi"
	"serviceapis/orderapi"
	"serviceapis/restaurantapi"
)

func RegisterTypes() {
	accountingapi.RegisterTypes()
	consumerapi.RegisterTypes()
	kitchenapi.RegisterTypes()
	orderapi.RegisterTypes()
	restaurantapi.RegisterTypes()
}
