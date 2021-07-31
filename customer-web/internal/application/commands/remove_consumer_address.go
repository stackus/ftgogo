package commands

import (
	"context"

	"github.com/stackus/ftgogo/customer-web/internal/adapters"
)

type RemoveConsumerAddress struct {
	ConsumerID string
	AddressID  string
}

type RemoveConsumerAddressHandler struct {
	repo adapters.ConsumerRepository
}

func NewRemoveConsumerAddressHandler(repo adapters.ConsumerRepository) RemoveConsumerAddressHandler {
	return RemoveConsumerAddressHandler{repo: repo}
}

func (h RemoveConsumerAddressHandler) Handle(ctx context.Context, cmd RemoveConsumerAddress) error {
	return h.repo.RemoveAddress(ctx, adapters.RemoveConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
	})
}
