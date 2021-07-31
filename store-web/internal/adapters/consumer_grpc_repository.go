package adapters

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/consumerapi/pb"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type ConsumerGRPCRepository struct {
	client consumerpb.ConsumerServiceClient
}

var _ ConsumerRepository = (*ConsumerGRPCRepository)(nil)

func NewConsumerGrpcRepository(client consumerpb.ConsumerServiceClient) *ConsumerGRPCRepository {
	return &ConsumerGRPCRepository{client: client}
}

func (r ConsumerGRPCRepository) Find(ctx context.Context, findConsumer FindConsumer) (*domain.Consumer, error) {
	resp, err := r.client.GetConsumer(ctx, &consumerpb.GetConsumerRequest{ConsumerID: findConsumer.ConsumerID})
	if err != nil {
		return nil, err
	}

	return &domain.Consumer{
		ConsumerID: resp.ConsumerID,
		Name:       resp.Name,
	}, nil
}
