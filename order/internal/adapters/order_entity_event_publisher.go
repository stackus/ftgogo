package adapters

import (
	"github.com/stackus/edat/msg"
)

func NewOrderEntityEventPublisher(publisher msg.EntityEventMessagePublisher) OrderPublisher {
	return publisher
}
