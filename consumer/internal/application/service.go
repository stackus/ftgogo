package application

import (
	"github.com/stackus/ftgogo/consumer/internal/application/commands"
	"github.com/stackus/ftgogo/consumer/internal/application/queries"
)

type Service struct {
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
