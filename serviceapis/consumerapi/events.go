package consumerapi

import (
	"github.com/stackus/edat/core"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

func registerEvents() {
	core.RegisterEvents(
		ConsumerRegistered{}, ConsumerUpdated{},
	)
}

type ConsumerEvent struct{}

func (ConsumerEvent) DestinationChannel() string { return ConsumerAggregateChannel }

type ConsumerRegistered struct {
	ConsumerEvent
	Name string
}

func (ConsumerRegistered) EventName() string { return "consumerapi.ConsumerRegistered" }

type ConsumerUpdated struct {
	ConsumerEvent
	Name string
}

func (ConsumerUpdated) EventName() string { return "consumerapi.ConsumerUpdated" }

type AddressAdded struct {
	ConsumerEvent
	AddressID string
	Address   *commonapi.Address
}

func (AddressAdded) EventName() string { return "consumerapi.AddressAdded" }

type AddressUpdated struct {
	ConsumerEvent
	AddressID string
	Address   *commonapi.Address
}

func (AddressUpdated) EventName() string { return "consumerapi.AddressUpdated" }

type AddressRemoved struct {
	ConsumerEvent
	AddressID string
}

func (AddressRemoved) EventName() string { return "consumerapi.AddressRemoved" }
