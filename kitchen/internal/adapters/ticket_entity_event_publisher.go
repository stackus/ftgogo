package adapters

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/kitchen/internal/application/ports"
)

func NewTicketEntityEventPublisher(publisher msg.EntityEventMessagePublisher) ports.TicketPublisher {
	return publisher
}
