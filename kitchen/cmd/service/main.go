package main

import (
	"github.com/stackus/ftgogo/kitchen/internal/adapters"
	"github.com/stackus/ftgogo/kitchen/internal/application"
	"github.com/stackus/ftgogo/kitchen/internal/domain"
	"github.com/stackus/ftgogo/kitchen/internal/handlers"
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
	domain.RegisterTypes()

	// Driven
	ticketRepo := adapters.NewTicketRepositoryPublisherMiddleware(
		adapters.NewTicketAggregateRepository(svc.AggregateStore),
		adapters.NewTicketEntityEventPublisher(svc.Publisher),
	)
	restaurantRepo := adapters.NewRestaurantPostgresRepository(svc.PgConn)

	app := application.NewServiceApplication(ticketRepo, restaurantRepo)

	// Drivers
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, svc.Publisher)
	handlers.NewRestaurantEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
