package queries

import (
	"context"

	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/application/ports"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type GetAddress struct {
	ConsumerID string
	AddressID  string
}

type GetAddressHandler struct {
	repo ports.ConsumerRepository
}

func NewGetAddressHandler(repo ports.ConsumerRepository) GetAddressHandler {
	return GetAddressHandler{repo: repo}
}

func (h GetAddressHandler) Handle(ctx context.Context, query GetAddress) (*commonapi.Address, error) {
	consumer, err := h.repo.Load(ctx, query.ConsumerID)
	if err != nil {
		return nil, err
	}

	address := consumer.Address(query.AddressID)
	if address == nil {
		return nil, errors.Wrap(errors.ErrNotFound, "an address with that identifier does not exist")
	}

	return address, nil
}
