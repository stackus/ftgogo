package commands

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type UpdateConsumerAddress struct {
	ConsumerID string
	AddressID  string
	Address    *commonapi.Address
}

type UpdateConsumerAddressHandler struct {
	repo domain.ConsumerRepository
}

func NewUpdateConsumerAddressHandler(repo domain.ConsumerRepository) UpdateConsumerAddressHandler {
	return UpdateConsumerAddressHandler{repo: repo}
}

func (h UpdateConsumerAddressHandler) Handle(ctx context.Context, cmd UpdateConsumerAddress) error {
	return h.repo.UpdateAddress(ctx, domain.UpdateConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
		Address:    cmd.Address,
	})
}
