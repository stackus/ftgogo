package adapters

import (
	"github.com/stackus/edat/es"
	"github.com/stackus/ftgogo/order/internal/domain"
)

func NewOrderRepository(store es.AggregateRootStore) domain.OrderRepository {
	return es.NewAggregateRootRepository(domain.NewOrder, store)
}
