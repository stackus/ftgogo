package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type AddAddress struct {
	ConsumerID string
	AddressID  string
	Address    *commonapi.Address
}

type AddAddressHandler struct {
	repo adapters.ConsumerRepository
}

func NewAddAddressHandler(repo adapters.ConsumerRepository) AddAddressHandler {
	return AddAddressHandler{repo: repo}
}

func (h AddAddressHandler) Handle(ctx context.Context, cmd AddAddress) error {
	_, err := h.repo.Update(ctx, cmd.ConsumerID, &domain.AddAddress{
		AddressID: cmd.AddressID,
		Address:   cmd.Address,
	})
	return err
}
