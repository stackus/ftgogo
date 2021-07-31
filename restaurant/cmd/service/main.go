package main

import (
	"github.com/stackus/ftgogo/restaurant/internal/adapters"
	"github.com/stackus/ftgogo/restaurant/internal/application"
	"github.com/stackus/ftgogo/restaurant/internal/handlers"
	"github.com/stackus/ftgogo/serviceapis"
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

	// Driven
	restaurantRepo := adapters.NewRestaurantPostgresPublisherMiddleware(
		adapters.NewRestaurantPostgresRepository(svc.PgConn),
		adapters.NewRestaurantEntityEventPublisher(svc.Publisher),
	)

	app := application.NewServiceApplication(restaurantRepo)

	// Drivers
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
