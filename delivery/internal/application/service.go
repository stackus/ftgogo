package application

import (
	"github.com/stackus/ftgogo/delivery/internal/application/commands"
	"github.com/stackus/ftgogo/delivery/internal/application/queries"
)

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateDelivery         commands.CreateDeliveryHandler
	CreateRestaurant       commands.CreateRestaurantHandler
	SetCourierAvailability commands.SetCourierAvailabilityHandler
	ScheduleDelivery       commands.ScheduleDeliveryHandler
	CancelDelivery         commands.CancelDeliveryHandler
}

type Queries struct {
	GetCourier  queries.GetCourierHandler
	GetDelivery queries.GetDeliveryHandler
}
