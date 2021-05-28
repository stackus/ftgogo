package application

import (
	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
	"github.com/stackus/ftgogo/restaurant/internal/application/queries"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateRestaurant commands.CreateRestaurantHandler
}

type Queries struct {
	GetRestaurant queries.GetRestaurantHandler
}
