package queries

import (
	"context"

	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/adapters"
	"github.com/stackus/ftgogo/consumer/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type GetAddress struct {
	ConsumerID string
	AddressID  string
}

type GetAddressHandler struct {
	repo adapters.ConsumerRepository
}

func NewGetAddressHandler(repo adapters.ConsumerRepository) GetAddressHandler {
	return GetAddressHandler{repo: repo}
}

func (h GetAddressHandler) Handle(ctx context.Context, query GetAddress) (*commonapi.Address, error) {
	root, err := h.repo.Load(ctx, query.ConsumerID)
	if err != nil {
		return nil, err
	}

	address := root.Aggregate().(*domain.Consumer).Address(query.AddressID)
	if address == nil {
		return nil, errors.Wrap(errors.ErrNotFound, "an address with that identifier does not exist")
	}

	return address, nil
}
