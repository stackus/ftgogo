package adapters

import (
	"github.com/stackus/edat/msg"
)

func NewConsumerEntityEventPublisher(publisher msg.EntityEventMessagePublisher) ConsumerPublisher {
	return publisher
}
