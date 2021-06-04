package adapters

import (
	"github.com/stackus/edat/msg"
)

func NewRestaurantEntityEventPublisher(publisher msg.EntityEventMessagePublisher) RestaurantPublisher {
	return publisher
}
