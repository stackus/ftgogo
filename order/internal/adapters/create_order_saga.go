package adapters

import (
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"github.com/stackus/ftgogo/order/internal/sagas"
)

func NewCreateOrderSaga(store saga.InstanceStore, publisher msg.CommandMessagePublisher, options ...saga.OrchestratorOption) *saga.Orchestrator {
	return saga.NewOrchestrator(sagas.NewCreateOrderSagaDefinition(), store, publisher, options...)
}
