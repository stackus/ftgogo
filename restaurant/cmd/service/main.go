package main

import (
	"github.com/stackus/ftgogo/restaurant/internal/adapters"
	"github.com/stackus/ftgogo/restaurant/internal/application"
	"github.com/stackus/ftgogo/restaurant/internal/application/commands"
	"github.com/stackus/ftgogo/restaurant/internal/application/queries"
	"github.com/stackus/ftgogo/restaurant/internal/ports"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/restaurantapi/pb"
	"shared-go/applications"
)

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

	app := application.Application{
		Commands: application.Commands{
			CreateRestaurant: commands.NewCreateRestaurantHandler(restaurantRepo),
		},
		Queries: application.Queries{
			GetRestaurant: queries.NewGetRestaurantHandler(restaurantRepo),
		},
	}

	restaurantpb.RegisterRestaurantServiceServer(svc.RpcServer, ports.NewRpcHandlers(app))

	return nil
}
