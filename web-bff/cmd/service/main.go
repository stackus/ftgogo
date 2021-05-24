package main

import (
	"context"

	"google.golang.org/grpc"

	consumerpb "github.com/stackus/ftgogo/serviceapis/consumerapi/pb"
	"github.com/stackus/ftgogo/serviceapis/orderapi/pb"
	"github.com/stackus/ftgogo/web-bff/internal/adapters"
	"github.com/stackus/ftgogo/web-bff/internal/application/commands"
	"github.com/stackus/ftgogo/web-bff/internal/application/queries"
	"shared-go/applications"
	"shared-go/rpc"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterConsumer      commands.RegisterConsumerHandler
	CreateOrder           commands.CreateOrderHandler
	AddConsumerAddress    commands.AddConsumerAddressHandler
	UpdateConsumerAddress commands.UpdateConsumerAddressHandler
	RemoveConsumerAddress commands.RemoveConsumerAddressHandler
}

type Queries struct {
	GetConsumer        queries.GetConsumerHandler
	GetOrder           queries.GetOrderHandler
	GetConsumerAddress queries.GetConsumerAddressHandler
}

func main() {
	svc := applications.NewService(initService)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initService(svc *applications.Service) error {
	var err error

	var orderConn *grpc.ClientConn
	orderConn, err = rpc.NewClientConn("order-service:8000", rpc.WithClientUnaryConvertStatus())
	if err != nil {
		return err
	}
	orderClient := adapters.NewOrderGrpcClient(orderpb.NewOrderServiceClient(orderConn))

	var consumerConn *grpc.ClientConn
	consumerConn, err = rpc.NewClientConn("consumer-service:8000", rpc.WithClientUnaryConvertStatus())
	if err != nil {
		return err
	}
	consumerClient := adapters.NewConsumerGrpcClient(consumerpb.NewConsumerServiceClient(consumerConn))

	application := Application{
		Commands: Commands{
			RegisterConsumer:      commands.NewRegisterConsumerHandler(consumerClient),
			CreateOrder:           commands.NewCreateOrderHandler(orderClient, consumerClient),
			AddConsumerAddress:    commands.NewAddConsumerAddressHandler(consumerClient),
			UpdateConsumerAddress: commands.NewUpdateConsumerAddressHandler(consumerClient),
			RemoveConsumerAddress: commands.NewRemoveConsumerAddressHandler(consumerClient),
		},
		Queries: Queries{
			GetConsumer:        queries.NewGetConsumerHandler(consumerClient),
			GetOrder:           queries.NewGetOrderHandler(orderClient),
			GetConsumerAddress: queries.NewGetConsumerAddressHandler(consumerClient),
		},
	}

	// api := newWebHandlers(application)
	svc.WebServer.Mount(svc.Cfg.Web.ApiPath, newWebHandlers(application).Mount)
	// svc.WebServer.Mount(svc.Cfg.Web.ApiPath, func(r chi.Router) http.Handler {
	// 	return HandlerFromMux(api, r)
	// })

	// Perform custom service cleanup
	svc.Cleanup = func(ctx context.Context) error {
		var err error

		// close all grpc connections
		err = orderConn.Close()

		return err
	}

	return nil
}
