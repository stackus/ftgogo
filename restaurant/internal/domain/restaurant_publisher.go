package domain

import (
	"github.com/stackus/edat/msg"
)

type RestaurantPublisher interface {
	msg.EntityEventMessagePublisher
}
