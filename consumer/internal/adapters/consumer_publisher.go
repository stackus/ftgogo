package adapters

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/consumer/internal/domain"
)

func NewConsumerPublisher(publisher msg.EntityEventMessagePublisher) domain.ConsumerPublisher {
	return publisher
}
