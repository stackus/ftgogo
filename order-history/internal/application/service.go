package application

import (
	"github.com/stackus/ftgogo/order-history/internal/application/commands"
	"github.com/stackus/ftgogo/order-history/internal/application/queries"
)

type Service struct {
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
