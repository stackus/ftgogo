package domain

import (
	"github.com/stackus/edat/core"
)

func registerConsumerCommands() {
	core.RegisterCommands(RegisterConsumer{})
}

type RegisterConsumer struct {
	Name string
}

func (RegisterConsumer) CommandName() string {
	return "consumerservice.RegisterConsumer"
}
