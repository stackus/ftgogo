package adapters

import (
	"github.com/stackus/edat/msg"

	"github.com/stackus/ftgogo/order/internal/domain"
)

func NewOrderPublisher(publisher msg.EntityEventMessagePublisher) domain.OrderPublisher {
	return publisher
}
