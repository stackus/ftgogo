package domain

import (
	"github.com/stackus/edat/core"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

func registerConsumerCommands() {
	core.RegisterCommands(
		RegisterConsumer{}, UpdateConsumer{},
		AddAddress{}, UpdateAddress{}, RemoveAddress{},
	)
}

type RegisterConsumer struct {
	Name string
}

func (RegisterConsumer) CommandName() string { return "consumerservice.RegisterConsumer" }

type UpdateConsumer struct {
	Name string
}

func (UpdateConsumer) CommandName() string { return "consumerservice.UpdateConsumer" }

type AddAddress struct {
	AddressID string
	Address   *commonapi.Address
}

func (AddAddress) CommandName() string { return "consumerservice.AddAddress" }

type UpdateAddress struct {
	AddressID string
	Address   *commonapi.Address
}

func (UpdateAddress) CommandName() string { return "consumerservice.UpdateAddress" }

type RemoveAddress struct {
	AddressID string
}

func (RemoveAddress) CommandName() string { return "consumerservice.RemoveAddress" }
