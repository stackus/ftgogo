package commands

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/application/ports"
	"github.com/stackus/ftgogo/consumer/internal/domain"
)

type RemoveAddress struct {
	ConsumerID string
	AddressID  string
}

type RemoveAddressHandler struct {
	repo ports.ConsumerRepository
}

func NewRemoveAddressHandler(repo ports.ConsumerRepository) RemoveAddressHandler {
	return RemoveAddressHandler{repo: repo}
}

func (h RemoveAddressHandler) Handle(ctx context.Context, cmd RemoveAddress) error {
	_, err := h.repo.Update(ctx, cmd.ConsumerID, &domain.RemoveAddress{
		AddressID: cmd.AddressID,
	})
	return err
}
