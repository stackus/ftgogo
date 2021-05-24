package adapters

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	commonpb "github.com/stackus/ftgogo/serviceapis/commonapi/pb"
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

func (r ConsumerGRPCRepository) Register(ctx context.Context, name string) (string, error) {
	resp, err := r.client.RegisterConsumer(ctx, &consumerpb.RegisterConsumerRequest{Name: name})
	if err != nil {
		return "", err
	}

	return resp.ConsumerID, nil
}

func (r ConsumerGRPCRepository) Find(ctx context.Context, consumerID string) (*domain.Consumer, error) {
	resp, err := r.client.GetConsumer(ctx, &consumerpb.GetConsumerRequest{ConsumerID: consumerID})
	if err != nil {
		return nil, err
	}

	return &domain.Consumer{
		ConsumerID: resp.ConsumerID,
		Name:       resp.Name,
	}, nil
}

func (r ConsumerGRPCRepository) Update(ctx context.Context, updateConsumer domain.Consumer) error {
	// _, err := r.client.UpdateConsumer(ctx, ...)
	panic("implement me")
}

// NOTE All of the address additions have been added to demonstrate BFF use cases (gather data from multiple services)

func (r ConsumerGRPCRepository) AddAddress(ctx context.Context, addAddress domain.ModifyConsumerAddress) error {
	_, err := r.client.AddAddress(ctx, &consumerpb.AddAddressRequest{
		ConsumerID: addAddress.ConsumerID,
		AddressID:  addAddress.AddressID,
		Address:    r.toAddressProto(addAddress.Address),
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
	return r.fromAddressProto(resp.Address), nil
}

func (r ConsumerGRPCRepository) UpdateAddress(ctx context.Context, updateAddress domain.ModifyConsumerAddress) error {
	_, err := r.client.UpdateAddress(ctx, &consumerpb.UpdateAddressRequest{
		ConsumerID: updateAddress.ConsumerID,
		AddressID:  updateAddress.AddressID,
		Address:    r.toAddressProto(updateAddress.Address),
	})
	return err
}

func (r ConsumerGRPCRepository) RemoveAddress(ctx context.Context, removeAddress domain.FindConsumerAddress) error {
	_, err := r.client.RemoveAddress(ctx, &consumerpb.RemoveAddressRequest{
		ConsumerID: removeAddress.ConsumerID,
		AddressID:  removeAddress.AddressID,
	})
	return err
}

func (r ConsumerGRPCRepository) toAddressProto(address *commonapi.Address) *commonpb.Address {
	return &commonpb.Address{
		Street1: address.Street1,
		Street2: address.Street2,
		City:    address.City,
		State:   address.State,
		Zip:     address.Zip,
	}
}

func (r ConsumerGRPCRepository) fromAddressProto(address *commonpb.Address) *commonapi.Address {
	return &commonapi.Address{
		Street1: address.Street1,
		Street2: address.Street2,
		City:    address.City,
		State:   address.State,
		Zip:     address.Zip,
	}
}
