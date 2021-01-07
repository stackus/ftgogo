package consumerapi

import (
	"github.com/stackus/edat/core"
)

func registerEvents() {
	core.RegisterEvents(ConsumerRegistered{})
}

type ConsumerEvent struct{}

func (ConsumerEvent) DestinationChannel() string { return ConsumerAggregateChannel }

type ConsumerRegistered struct {
	ConsumerEvent
	Name string
}

func (ConsumerRegistered) EventName() string { return "consumerapi.ConsumerRegistered" }
