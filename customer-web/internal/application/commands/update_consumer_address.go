package commands

import (
	"context"

	"github.com/stackus/ftgogo/customer-web/internal/adapters"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type UpdateConsumerAddress struct {
	ConsumerID string
	AddressID  string
	Address    *commonapi.Address
}

type UpdateConsumerAddressHandler struct {
	repo adapters.ConsumerRepository
}

func NewUpdateConsumerAddressHandler(repo adapters.ConsumerRepository) UpdateConsumerAddressHandler {
	return UpdateConsumerAddressHandler{repo: repo}
}

func (h UpdateConsumerAddressHandler) Handle(ctx context.Context, cmd UpdateConsumerAddress) error {
	return h.repo.UpdateAddress(ctx, adapters.UpdateConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
		Address:    cmd.Address,
	})
}
