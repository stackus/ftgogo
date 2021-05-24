package commands

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type AddConsumerAddress struct {
	ConsumerID string
	AddressID  string
	Address    *commonapi.Address
}

type AddConsumerAddressHandler struct {
	repo domain.ConsumerRepository
}

func NewAddConsumerAddressHandler(repo domain.ConsumerRepository) AddConsumerAddressHandler {
	return AddConsumerAddressHandler{repo: repo}
}

func (h AddConsumerAddressHandler) Handle(ctx context.Context, cmd AddConsumerAddress) error {
	return h.repo.AddAddress(ctx, domain.ModifyConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
		Address:    cmd.Address,
	})
}
