package application

import (
	"context"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
	"github.com/stackus/ftgogo/consumer/internal/application/ports"
	"github.com/stackus/ftgogo/consumer/internal/application/queries"
	"github.com/stackus/ftgogo/consumer/internal/domain"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type ServiceApplication interface {
	RegisterConsumer(ctx context.Context, cmd commands.RegisterConsumer) (string, error)
	UpdateConsumer(ctx context.Context, cmd commands.UpdateConsumer) error
	ValidateOrderByConsumer(ctx context.Context, cmd commands.ValidateOrderByConsumer) error
	AddAddress(ctx context.Context, cmd commands.AddAddress) error
	UpdateAddress(ctx context.Context, cmd commands.UpdateAddress) error
	RemoveAddress(ctx context.Context, cmd commands.RemoveAddress) error
	GetConsumer(ctx context.Context, query queries.GetConsumer) (*domain.Consumer, error)
	GetAddress(ctx context.Context, query queries.GetAddress) (*commonapi.Address, error)
}

type Service struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterConsumer        commands.RegisterConsumerHandler
	UpdateConsumer          commands.UpdateConsumerHandler
	ValidateOrderByConsumer commands.ValidateOrderByConsumerHandler
	AddAddress              commands.AddAddressHandler
	UpdateAddress           commands.UpdateAddressHandler
	RemoveAddress           commands.RemoveAddressHandler
}

type Queries struct {
	GetConsumer queries.GetConsumerHandler
	GetAddress  queries.GetAddressHandler
}

func NewServiceApplication(consumerRepo ports.ConsumerRepository) *Service {
	return &Service{
		Commands: Commands{
			RegisterConsumer:        commands.NewRegisterConsumerHandler(consumerRepo),
			UpdateConsumer:          commands.NewUpdateConsumerHandler(consumerRepo),
			ValidateOrderByConsumer: commands.NewValidateOrderByConsumerHandler(consumerRepo),
			AddAddress:              commands.NewAddAddressHandler(consumerRepo),
			UpdateAddress:           commands.NewUpdateAddressHandler(consumerRepo),
			RemoveAddress:           commands.NewRemoveAddressHandler(consumerRepo),
		},
		Queries: Queries{
			GetConsumer: queries.NewGetConsumerHandler(consumerRepo),
			GetAddress:  queries.NewGetAddressHandler(consumerRepo),
		},
	}
}

func (s Service) RegisterConsumer(ctx context.Context, cmd commands.RegisterConsumer) (string, error) {
	return s.Commands.RegisterConsumer.Handle(ctx, cmd)
}

func (s Service) UpdateConsumer(ctx context.Context, cmd commands.UpdateConsumer) error {
	return s.Commands.UpdateConsumer.Handle(ctx, cmd)
}

func (s Service) ValidateOrderByConsumer(ctx context.Context, cmd commands.ValidateOrderByConsumer) error {
	return s.Commands.ValidateOrderByConsumer.Handle(ctx, cmd)
}

func (s Service) AddAddress(ctx context.Context, cmd commands.AddAddress) error {
	return s.Commands.AddAddress.Handle(ctx, cmd)
}

func (s Service) UpdateAddress(ctx context.Context, cmd commands.UpdateAddress) error {
	return s.Commands.UpdateAddress.Handle(ctx, cmd)
}

func (s Service) RemoveAddress(ctx context.Context, cmd commands.RemoveAddress) error {
	return s.Commands.RemoveAddress.Handle(ctx, cmd)
}

func (s Service) GetConsumer(ctx context.Context, query queries.GetConsumer) (*domain.Consumer, error) {
	return s.Queries.GetConsumer.Handle(ctx, query)
}

func (s Service) GetAddress(ctx context.Context, query queries.GetAddress) (*commonapi.Address, error) {
	return s.Queries.GetAddress.Handle(ctx, query)
}
