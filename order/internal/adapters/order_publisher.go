package adapters

import (
	"github.com/stackus/edat/msg"
)

type OrderPublisher interface {
	msg.EntityEventMessagePublisher
}
