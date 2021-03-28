package adapters

import (
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/consumer/internal/domain"
)

func NewConsumerRepository(store es.AggregateRootStore) domain.ConsumerRepository {
	return es.NewAggregateRootRepository(domain.NewConsumer, store)
}
