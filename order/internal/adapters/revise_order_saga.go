package adapters

import (
	"context"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"
)

type ReviseOrderSaga interface {
	Start(ctx context.Context, sagaData core.SagaData) (*saga.Instance, error)
	ReplyChannel() string
	ReceiveMessage(ctx context.Context, message msg.Message) error
}
