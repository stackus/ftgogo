package adapters

import (
	"github.com/stackus/edat/es"

	"github.com/stackus/ftgogo/accounting/internal/domain"
)

func NewAccountAggregateRootRepository(store es.AggregateRootStore) AccountRepository {
	return es.NewAggregateRootRepository(domain.NewAccount, store)
}
