package adapters

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/kitchen/internal/domain"
)

func NewTicketPublisher(publisher msg.EntityEventMessagePublisher) domain.TicketPublisher {
	return publisher
}
