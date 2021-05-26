package commands

import (
	"context"

	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type RemoveConsumerAddress struct {
	ConsumerID string
	AddressID  string
}

type RemoveConsumerAddressHandler struct {
	repo domain.ConsumerRepository
}

func NewRemoveConsumerAddressHandler(repo domain.ConsumerRepository) RemoveConsumerAddressHandler {
	return RemoveConsumerAddressHandler{repo: repo}
}

func (h RemoveConsumerAddressHandler) Handle(ctx context.Context, cmd RemoveConsumerAddress) error {
	return h.repo.RemoveAddress(ctx, domain.RemoveConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
	})
}
