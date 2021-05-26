package adapters

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	consumerpb "github.com/stackus/ftgogo/serviceapis/consumerapi/pb"
	"github.com/stackus/ftgogo/web-bff/internal/domain"
)

type ConsumerGRPCRepository struct {
	client consumerpb.ConsumerServiceClient
}

var _ domain.ConsumerRepository = (*ConsumerGRPCRepository)(nil)

func NewConsumerGrpcClient(client consumerpb.ConsumerServiceClient) *ConsumerGRPCRepository {
	return &ConsumerGRPCRepository{client: client}
}

func (r ConsumerGRPCRepository) Register(ctx context.Context, registerConsumer domain.RegisterConsumer) (string, error) {
	resp, err := r.client.RegisterConsumer(ctx, &consumerpb.RegisterConsumerRequest{Name: registerConsumer.Name})
	if err != nil {
		return "", err
	}

	return resp.ConsumerID, nil
}

func (r ConsumerGRPCRepository) Find(ctx context.Context, findConsumer domain.FindConsumer) (*domain.Consumer, error) {
	resp, err := r.client.GetConsumer(ctx, &consumerpb.GetConsumerRequest{ConsumerID: findConsumer.ConsumerID})
	if err != nil {
		return nil, err
	}

	return &domain.Consumer{
		ConsumerID: resp.ConsumerID,
		Name:       resp.Name,
	}, nil
}

func (r ConsumerGRPCRepository) Update(ctx context.Context, updateConsumer domain.UpdateConsumer) error {
	// _, err := r.client.UpdateConsumer(ctx, ...)
	panic("implement me")
}

// NOTE All of the address additions have been added to demonstrate BFF use cases (gather data from multiple services)

func (r ConsumerGRPCRepository) AddAddress(ctx context.Context, addAddress domain.AddConsumerAddress) error {
	_, err := r.client.AddAddress(ctx, &consumerpb.AddAddressRequest{
		ConsumerID: addAddress.ConsumerID,
		AddressID:  addAddress.AddressID,
		Address:    commonapi.ToAddressProto(addAddress.Address),
	})
	return err
}

func (r ConsumerGRPCRepository) FindAddress(ctx context.Context, findAddress domain.FindConsumerAddress) (*commonapi.Address, error) {
	resp, err := r.client.GetAddress(ctx, &consumerpb.GetAddressRequest{
		ConsumerID: findAddress.ConsumerID,
		AddressID:  findAddress.AddressID,
	})
	if err != nil {
		return nil, err
	}
	return commonapi.FromAddressProto(resp.Address), nil
}

func (r ConsumerGRPCRepository) UpdateAddress(ctx context.Context, updateAddress domain.UpdateConsumerAddress) error {
	_, err := r.client.UpdateAddress(ctx, &consumerpb.UpdateAddressRequest{
		ConsumerID: updateAddress.ConsumerID,
		AddressID:  updateAddress.AddressID,
		Address:    commonapi.ToAddressProto(updateAddress.Address),
	})
	return err
}

func (r ConsumerGRPCRepository) RemoveAddress(ctx context.Context, removeAddress domain.RemoveConsumerAddress) error {
	_, err := r.client.RemoveAddress(ctx, &consumerpb.RemoveAddressRequest{
		ConsumerID: removeAddress.ConsumerID,
		AddressID:  removeAddress.AddressID,
	})
	return err
}
