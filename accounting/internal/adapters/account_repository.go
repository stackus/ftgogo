package adapters

import (
	"github.com/stackus/edat/es"
	"github.com/stackus/ftgogo/account/internal/domain"
)

func NewAccountRepository(store es.AggregateRootStore) domain.AccountRepository {
	return es.NewAggregateRootRepository(domain.NewAccount, store)
}
