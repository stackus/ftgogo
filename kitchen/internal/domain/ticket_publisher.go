package domain

import (
	"github.com/stackus/edat/msg"
)

type TicketPublisher interface {
	msg.EntityEventMessagePublisher
}
