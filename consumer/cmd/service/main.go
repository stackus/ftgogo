package main

import (
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/application/commands"
	"github.com/stackus/ftgogo/consumer/internal/application/queries"
	"github.com/stackus/ftgogo/consumer/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/consumerapi"
	"github.com/stackus/ftgogo/serviceapis/consumerapi/pb"
	"shared-go/applications"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterConsumer        commands.RegisterConsumerHandler
	UpdateConsumer          commands.UpdateConsumerHandler
	ValidateOrderByConsumer commands.ValidateOrderByConsumerHandler
	AddAddress              commands.AddAddressHandler
	UpdateAddress           commands.UpdateAddressHandler
	RemoveAddress           commands.RemoveAddressHandler
}

type Queries struct {
	GetConsumer queries.GetConsumerHandler
	GetAddress  queries.GetAddressHandler
}

func main() {
	svc := applications.NewService(initService)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initService(svc *applications.Service) error {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()

	consumerRepo := adapters.NewConsumerRepositoryPublisherMiddleware(
		adapters.NewConsumerRepository(svc.AggregateStore),
		adapters.NewConsumerPublisher(svc.Publisher),
	)

	application := Application{
		Commands: Commands{
			RegisterConsumer:        commands.NewRegisterConsumerHandler(consumerRepo),
			UpdateConsumer:          commands.NewUpdateConsumerHandler(consumerRepo),
			ValidateOrderByConsumer: commands.NewValidateOrderByConsumerHandler(consumerRepo),
			AddAddress:              commands.NewAddAddressHandler(consumerRepo),
			UpdateAddress:           commands.NewUpdateAddressHandler(consumerRepo),
			RemoveAddress:           commands.NewRemoveAddressHandler(consumerRepo),
		},
		Queries: Queries{
			GetConsumer: queries.NewGetConsumerHandler(consumerRepo),
			GetAddress:  queries.NewGetAddressHandler(consumerRepo),
		},
	}

	cmdHandlers := NewCommandHandlers(application)
	svc.Subscriber.Subscribe(consumerapi.ConsumerServiceCommandChannel, saga.NewCommandDispatcher(svc.Publisher).
		Handle(consumerapi.ValidateOrderByConsumer{}, cmdHandlers.ValidateOrderByConsumer))

	consumerpb.RegisterConsumerServiceServer(svc.RpcServer, newRpcHandlers(application))

	return nil
}
