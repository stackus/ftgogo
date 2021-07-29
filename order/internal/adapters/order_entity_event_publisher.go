package adapters

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/order/internal/application/ports"
)

func NewOrderEntityEventPublisher(publisher msg.EntityEventMessagePublisher) ports.OrderPublisher {
	return publisher
}
