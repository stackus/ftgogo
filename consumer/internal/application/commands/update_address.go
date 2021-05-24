package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type UpdateAddress struct {
	ConsumerID string
	AddressID  string
	Address    *commonapi.Address
}

type UpdateAddressHandler struct {
	repo domain.ConsumerRepository
}

func NewUpdateAddressHandler(repo domain.ConsumerRepository) UpdateAddressHandler {
	return UpdateAddressHandler{repo: repo}
}

func (h UpdateAddressHandler) Handle(ctx context.Context, cmd UpdateAddress) error {
	_, err := h.repo.Update(ctx, cmd.ConsumerID, &domain.UpdateAddress{
		AddressID: cmd.AddressID,
		Address:   cmd.Address,
	})
	return err
}
