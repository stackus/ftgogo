package adapters

import (
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/order/internal/domain"
)

func NewOrderAggregateRootRepository(store es.AggregateRootStore) OrderRepository {
	return es.NewAggregateRootRepository(domain.NewOrder, store)
}
