package main

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/order-history/internal/adapters"
	"github.com/stackus/ftgogo/order-history/internal/application/commands"
	"github.com/stackus/ftgogo/order-history/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/orderapi"
	"github.com/stackus/ftgogo/serviceapis/orderhistoryapi/pb"
	"shared-go/applications"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrderHistory commands.CreateOrderHistoryHandler
	UpdateOrderStatus  commands.UpdateOrderStatusHandler
}

type Queries struct {
	SearchOrderHistories queries.SearchOrderHistoriesHandler
	GetOrderHistory      queries.GetOrderHistoryHandler
}

func main() {
	svc := applications.NewService(initService)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initService(svc *applications.Service) error {
	serviceapis.RegisterTypes()

	orderHistoryRepo := adapters.NewOrderHistoryPostgresRepository(svc.PgConn)

	application := Application{
		Commands: Commands{
			CreateOrderHistory: commands.NewCreateOrderHistoryHandler(orderHistoryRepo),
			UpdateOrderStatus:  commands.NewUpdateOrderStatusHandler(orderHistoryRepo),
		},
		Queries: Queries{
			SearchOrderHistories: queries.NewSearchOrderHistoriesHandler(orderHistoryRepo),
			GetOrderHistory:      queries.NewGetOrderHistoryHandler(orderHistoryRepo),
		},
	}

	orderEventHandlers := NewOrderEventHandlers(application)
	svc.Subscriber.Subscribe(orderapi.OrderAggregateChannel, msg.NewEntityEventDispatcher().
		Handle(orderapi.OrderCreated{}, orderEventHandlers.OrderCreated).
		Handle(orderapi.OrderApproved{}, orderEventHandlers.OrderApproved).
		Handle(orderapi.OrderCancelled{}, orderEventHandlers.OrderCancelled).
		Handle(orderapi.OrderRejected{}, orderEventHandlers.OrderRejected))

	orderhistorypb.RegisterOrderHistoryServiceServer(svc.RpcServer, newRpcHandlers(application))

	return nil
}
