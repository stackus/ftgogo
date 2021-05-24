package queries

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type GetConsumerAddress struct {
	ConsumerID string
	AddressID  string
}

type GetConsumerAddressHandler struct {
	repo domain.ConsumerRepository
}

func NewGetConsumerAddressHandler(repo domain.ConsumerRepository) GetConsumerAddressHandler {
	return GetConsumerAddressHandler{repo: repo}
}

func (h GetConsumerAddressHandler) Handle(ctx context.Context, cmd GetConsumerAddress) (*commonapi.Address, error) {
	return h.repo.FindAddress(ctx, domain.FindConsumerAddress{
		ConsumerID: cmd.ConsumerID,
		AddressID:  cmd.AddressID,
	})
}
