package adapters

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/restaurant/internal/application/ports"
)

func NewRestaurantEntityEventPublisher(publisher msg.EntityEventMessagePublisher) ports.RestaurantPublisher {
	return publisher
}
