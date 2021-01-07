package adapters

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/restaurant/internal/domain"
)

func NewRestaurantPublisher(publisher msg.EntityEventMessagePublisher) domain.RestaurantPublisher {
	return publisher
}
