package application

import (
	"github.com/stackus/ftgogo/customer-web/internal/application/commands"
	"github.com/stackus/ftgogo/customer-web/internal/application/queries"
)

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterConsumer      commands.RegisterConsumerHandler
	CreateOrder           commands.CreateOrderHandler
	ReviseOrder           commands.ReviseOrderHandler
	CancelOrder           commands.CancelOrderHandler
	AddConsumerAddress    commands.AddConsumerAddressHandler
	UpdateConsumerAddress commands.UpdateConsumerAddressHandler
	RemoveConsumerAddress commands.RemoveConsumerAddressHandler
}

type Queries struct {
	GetConsumer        queries.GetConsumerHandler
	GetOrder           queries.GetOrderHandler
	GetConsumerAddress queries.GetConsumerAddressHandler
	SearchOrders       queries.SearchOrdersHandler
}
