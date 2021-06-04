package adapters

import (
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/consumer/internal/domain"
)

func NewConsumerAggregateRootRepository(store es.AggregateRootStore) ConsumerRepository {
	return es.NewAggregateRootRepository(domain.NewConsumer, store)
}
