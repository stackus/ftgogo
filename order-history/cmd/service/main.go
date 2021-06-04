package main

import (
	"github.com/stackus/ftgogo/order-history/internal/adapters"
	"github.com/stackus/ftgogo/order-history/internal/application"
	"github.com/stackus/ftgogo/order-history/internal/application/commands"
	"github.com/stackus/ftgogo/order-history/internal/application/queries"
	"github.com/stackus/ftgogo/order-history/internal/handlers"
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
	orderHistoryRepo := adapters.NewOrderHistoryPostgresRepository(svc.PgConn)

	app := application.Service{
		Commands: application.Commands{
			CreateOrderHistory: commands.NewCreateOrderHistoryHandler(orderHistoryRepo),
			UpdateOrderStatus:  commands.NewUpdateOrderStatusHandler(orderHistoryRepo),
		},
		Queries: application.Queries{
			SearchOrderHistories: queries.NewSearchOrderHistoriesHandler(orderHistoryRepo),
			GetOrderHistory:      queries.NewGetOrderHistoryHandler(orderHistoryRepo),
		},
	}

	// Drivers
	handlers.NewOrderEventHandlers(app).Mount(svc.Subscriber)
	handlers.NewRpcHandlers(app).Mount(svc.RpcServer)

	return nil
}
