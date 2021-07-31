package serviceapis

import (
	"github.com/stackus/ftgogo/serviceapis/accountingapi"
	"github.com/stackus/ftgogo/serviceapis/consumerapi"
	"github.com/stackus/ftgogo/serviceapis/kitchenapi"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi"
)

func RegisterTypes() {
	accountingapi.RegisterTypes()
	consumerapi.RegisterTypes()
	kitchenapi.RegisterTypes()
	orderapi.RegisterTypes()
	restaurantapi.RegisterTypes()
}
