package consumerapi

import (
	"github.com/stackus/edat/core"
)

func registerCommands() {
	core.RegisterCommands(ValidateOrderByConsumer{})
}

type ConsumerServiceCommand struct{}

func (ConsumerServiceCommand) DestinationChannel() string { return ConsumerServiceCommandChannel }

type ValidateOrderByConsumer struct {
	ConsumerServiceCommand
	ConsumerID string
	OrderID    string
	OrderTotal int // Money
}

func (ValidateOrderByConsumer) CommandName() string { return "consumerapi.ValidateOrderByConsumer" }
