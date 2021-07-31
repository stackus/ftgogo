package ohismod

import (
	"github.com/stackus/ftgogo/order-history/internal/adapters"
	"github.com/stackus/ftgogo/order-history/internal/application"
	"github.com/stackus/ftgogo/order-history/internal/application/commands"
	"github.com/stackus/ftgogo/order-history/internal/application/queries"
	"github.com/stackus/ftgogo/order-history/internal/handlers"
	"shared-go/applications"
)

func Setup(svc *applications.Monolith) error {
	// Driven
	adapters.OrderHistoriesTableName = "orderhistory.orders"
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
