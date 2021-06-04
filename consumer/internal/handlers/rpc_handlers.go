package handlers

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stackus/ftgogo/consumer/internal/application"
	"github.com/stackus/ftgogo/consumer/internal/application/commands"
	"github.com/stackus/ftgogo/consumer/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/consumerapi/pb"
)

type RpcHandlers struct {
	app application.Service
	consumerpb.UnimplementedConsumerServiceServer
}

var _ consumerpb.ConsumerServiceServer = (*RpcHandlers)(nil)

func NewRpcHandlers(app application.Service) RpcHandlers {
	return RpcHandlers{app: app}
}

func (h RpcHandlers) Mount(registrar grpc.ServiceRegistrar) {
	consumerpb.RegisterConsumerServiceServer(registrar, h)
}

func (h RpcHandlers) RegisterConsumer(ctx context.Context, request *consumerpb.RegisterConsumerRequest) (*consumerpb.RegisterConsumerResponse, error) {
	consumerID, err := h.app.Commands.RegisterConsumer.Handle(ctx, commands.RegisterConsumer{Name: request.Name})
	if err != nil {
		return nil, err
	}
	return &consumerpb.RegisterConsumerResponse{ConsumerID: consumerID}, nil
}

func (h RpcHandlers) GetConsumer(ctx context.Context, request *consumerpb.GetConsumerRequest) (*consumerpb.GetConsumerResponse, error) {
	consumer, err := h.app.Queries.GetConsumer.Handle(ctx, queries.GetConsumer{ConsumerID: request.ConsumerID})
	if err != nil {
		return nil, err
	}
	return &consumerpb.GetConsumerResponse{
		ConsumerID: consumer.ID(),
		Name:       consumer.Name(),
	}, nil
}

func (h RpcHandlers) UpdateConsumer(ctx context.Context, request *consumerpb.UpdateConsumerRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.UpdateConsumer.Handle(ctx, commands.UpdateConsumer{
		ConsumerID: request.ConsumerID,
		Name:       request.Name,
	})
	return &emptypb.Empty{}, err
}

func (h RpcHandlers) AddAddress(ctx context.Context, request *consumerpb.AddAddressRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.AddAddress.Handle(ctx, commands.AddAddress{
		ConsumerID: request.ConsumerID,
		AddressID:  request.AddressID,
		Address:    commonapi.FromAddressProto(request.Address),
	})
	return &emptypb.Empty{}, err
}

func (h RpcHandlers) GetAddress(ctx context.Context, request *consumerpb.GetAddressRequest) (*consumerpb.GetAddressResponse, error) {
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
		Address:    commonapi.ToAddressProto(address),
	}, nil
}

func (h RpcHandlers) UpdateAddress(ctx context.Context, request *consumerpb.UpdateAddressRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.UpdateAddress.Handle(ctx, commands.UpdateAddress{
		ConsumerID: request.ConsumerID,
		AddressID:  request.AddressID,
		Address:    commonapi.FromAddressProto(request.Address),
	})
	return &emptypb.Empty{}, err
}

func (h RpcHandlers) RemoveAddress(ctx context.Context, request *consumerpb.RemoveAddressRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.RemoveAddress.Handle(ctx, commands.RemoveAddress{
		ConsumerID: request.ConsumerID,
		AddressID:  request.AddressID,
	})
	return &emptypb.Empty{}, err
}
