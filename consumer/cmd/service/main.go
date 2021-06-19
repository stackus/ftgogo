package main

import (
	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/application"
	"github.com/stackus/ftgogo/consumer/internal/domain"
	"github.com/stackus/ftgogo/consumer/internal/handlers"
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
	consumerRepo := adapters.NewConsumerRepositoryPublisherMiddleware(
		adapters.NewConsumerAggregateRepository(svc.AggregateStore),
		adapters.NewConsumerEntityEventPublisher(svc.Publisher),
	)

	app := application.NewServiceApplication(consumerRepo)

	// Drivers
	handlers.NewCommandHandlers(app).Mount(svc.Subscriber, svc.Publisher)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
