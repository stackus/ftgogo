package adapters

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/consumer/internal/application/ports"
)

func NewConsumerEntityEventPublisher(publisher msg.EntityEventMessagePublisher) ports.ConsumerPublisher {
	return publisher
}
