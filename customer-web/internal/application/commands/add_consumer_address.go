package commands

import (
	"context"

	"github.com/stackus/ftgogo/customer-web/internal/adapters"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type AddConsumerAddress struct {
	ConsumerID string
	AddressID  string
	Address    *commonapi.Address
}

type AddConsumerAddressHandler struct {
	repo adapters.ConsumerRepository
}

func NewAddConsumerAddressHandler(repo adapters.ConsumerRepository) AddConsumerAddressHandler {
	return AddConsumerAddressHandler{repo: repo}
}

func (h AddConsumerAddressHandler) Handle(ctx context.Context, cmd AddConsumerAddress) error {
	return h.repo.AddAddress(ctx, adapters.AddConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
		Address:    cmd.Address,
	})
}
