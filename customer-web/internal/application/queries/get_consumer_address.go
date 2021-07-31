package queries

import (
	"context"

	"github.com/stackus/ftgogo/customer-web/internal/adapters"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type GetConsumerAddress struct {
	ConsumerID string
	AddressID  string
}

type GetConsumerAddressHandler struct {
	repo adapters.ConsumerRepository
}

func NewGetConsumerAddressHandler(repo adapters.ConsumerRepository) GetConsumerAddressHandler {
	return GetConsumerAddressHandler{repo: repo}
}

func (h GetConsumerAddressHandler) Handle(ctx context.Context, cmd GetConsumerAddress) (*commonapi.Address, error) {
	return h.repo.FindAddress(ctx, adapters.FindConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
	})
}
