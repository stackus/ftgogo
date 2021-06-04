package adapters

import (
	"github.com/stackus/edat/msg"
)

type ConsumerPublisher interface {
	msg.EntityEventMessagePublisher
}
