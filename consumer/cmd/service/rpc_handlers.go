package main

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
	"github.com/stackus/ftgogo/consumer/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	commonpb "github.com/stackus/ftgogo/serviceapis/commonapi/pb"
	"github.com/stackus/ftgogo/serviceapis/consumerapi/pb"
)

type rpcHandlers struct {
	app Application
	consumerpb.UnimplementedConsumerServiceServer
}

var _ consumerpb.ConsumerServiceServer = (*rpcHandlers)(nil)

func newRpcHandlers(app Application) rpcHandlers {
	return rpcHandlers{app: app}
}

func (h rpcHandlers) RegisterConsumer(ctx context.Context, request *consumerpb.RegisterConsumerRequest) (*consumerpb.RegisterConsumerResponse, error) {
	consumerID, err := h.app.Commands.RegisterConsumer.Handle(ctx, commands.RegisterConsumer{Name: request.Name})
	if err != nil {
		return nil, err
	}
	return &consumerpb.RegisterConsumerResponse{ConsumerID: consumerID}, nil
}

func (h rpcHandlers) GetConsumer(ctx context.Context, request *consumerpb.GetConsumerRequest) (*consumerpb.GetConsumerResponse, error) {
	consumer, err := h.app.Queries.GetConsumer.Handle(ctx, queries.GetConsumer{ConsumerID: request.ConsumerID})
	if err != nil {
		return nil, err
	}
	return &consumerpb.GetConsumerResponse{
		ConsumerID: consumer.ID(),
		Name:       consumer.Name(),
	}, nil
}

func (h rpcHandlers) UpdateConsumer(ctx context.Context, request *consumerpb.UpdateConsumerRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.UpdateConsumer.Handle(ctx, commands.UpdateConsumer{
		ConsumerID: request.ConsumerID,
		Name:       request.Name,
	})
	return &emptypb.Empty{}, err
}

func (h rpcHandlers) AddAddress(ctx context.Context, request *consumerpb.AddAddressRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.AddAddress.Handle(ctx, commands.AddAddress{
		ConsumerID: request.ConsumerID,
		AddressID:  request.AddressID,
		Address:    h.fromAddressProto(request.Address),
	})
	return &emptypb.Empty{}, err
}

func (h rpcHandlers) GetAddress(ctx context.Context, request *consumerpb.GetAddressRequest) (*consumerpb.GetAddressResponse, error) {
	address, err := h.app.Queries.GetAddress.Handle(ctx, queries.GetAddress{
		ConsumerID: request.ConsumerID,
		AddressID:  request.AddressID,
	})
	if err != nil {
		return nil, err
	}

	return &consumerpb.GetAddressResponse{
		ConsumerID: request.ConsumerID,
		AddressID:  request.AddressID,
		Address:    h.toAddressProto(address),
	}, nil
}

func (h rpcHandlers) UpdateAddress(ctx context.Context, request *consumerpb.UpdateAddressRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.UpdateAddress.Handle(ctx, commands.UpdateAddress{
		ConsumerID: request.ConsumerID,
		AddressID:  request.AddressID,
		Address:    h.fromAddressProto(request.Address),
	})
	return &emptypb.Empty{}, err
}

func (h rpcHandlers) RemoveAddress(ctx context.Context, request *consumerpb.RemoveAddressRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.RemoveAddress.Handle(ctx, commands.RemoveAddress{
		ConsumerID: request.ConsumerID,
		AddressID:  request.AddressID,
	})
	return &emptypb.Empty{}, err
}

func (h rpcHandlers) toAddressProto(address *commonapi.Address) *commonpb.Address {
	return &commonpb.Address{
		Street1: address.Street1,
		Street2: address.Street2,
		City:    address.City,
		State:   address.State,
		Zip:     address.Zip,
	}
}

func (h rpcHandlers) fromAddressProto(address *commonpb.Address) *commonapi.Address {
	return &commonapi.Address{
		Street1: address.Street1,
		Street2: address.Street2,
		City:    address.City,
		State:   address.State,
		Zip:     address.Zip,
	}
}
