package adapters

import (
	"github.com/stackus/edat/msg"
)

func NewTicketEntityEventPublisher(publisher msg.EntityEventMessagePublisher) TicketPublisher {
	return publisher
}
