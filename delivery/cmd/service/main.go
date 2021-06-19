package main

import (
	"github.com/stackus/ftgogo/delivery/internal/adapters"
	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/handlers"
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
	courierRepo := adapters.NewCourierPostgresRepository(svc.PgConn)
	deliveryRepo := adapters.NewDeliveryPostgresRepository(svc.PgConn)
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	app := application.NewServiceApplication(courierRepo, deliveryRepo, restaurantRepo)

	// Drivers
	handlers.NewRestaurantEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewOrderEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewTicketEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
