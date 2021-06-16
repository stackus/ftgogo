package application

import (
	"github.com/stackus/ftgogo/store-web/internal/application/commands"
	"github.com/stackus/ftgogo/store-web/internal/application/queries"
)

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	EnableAccount          commands.EnableAccountHandler
	DisableAccount         commands.DisableAccountHandler
	SetCourierAvailability commands.SetCourierAvailabilityHandler
	CancelOrder            commands.CancelOrderHandler
	CreateRestaurant       commands.CreateRestaurantHandler
}

type Queries struct {
	GetAccount         queries.GetAccountHandler
	GetConsumer        queries.GetConsumerHandler
	GetDeliveryHistory queries.GetDeliveryHistoryHandler
	GetOrder           queries.GetOrderHandler
	GetRestaurant      queries.GetRestaurantHandler
}
