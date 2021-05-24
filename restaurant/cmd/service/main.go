package main

import (
	"github.com/stackus/ftgogo/restaurant/internal/adapters"
	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
	"github.com/stackus/ftgogo/restaurant/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi/pb"
	"shared-go/applications"
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

func main() {
	svc := applications.NewService(initService)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initService(svc *applications.Service) error {
	serviceapis.RegisterTypes()

	restaurantRepo := adapters.NewRestaurantPostgresPublisherMiddleware(
		adapters.NewRestaurantPostgresRepository(svc.PgConn),
		adapters.NewRestaurantPublisher(svc.Publisher),
	)

	application := Application{
		Commands: Commands{
			CreateRestaurant: commands.NewCreateRestaurantHandler(restaurantRepo),
		},
		Queries: Queries{
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}

	restaurantpb.RegisterRestaurantServiceServer(svc.RpcServer, newRpcHandlers(application))

	return nil
}
