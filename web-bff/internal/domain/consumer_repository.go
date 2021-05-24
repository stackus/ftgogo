package domain

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type Consumer struct {
	ConsumerID string
	Name       string
	// TODO Addresses map[string]consumerapi.Address ??
}

type (
	ModifyConsumerAddress struct {
		ConsumerID string
		AddressID  string
		Address    *commonapi.Address
	}

	FindConsumerAddress struct {
		ConsumerID string
		AddressID  string
	}
)

type ConsumerRepository interface {
	Register(ctx context.Context, name string) (string, error)
	Find(ctx context.Context, consumerID string) (*Consumer, error)
	Update(ctx context.Context, updateConsumer Consumer) error
	AddAddress(ctx context.Context, addAddress ModifyConsumerAddress) error
	FindAddress(ctx context.Context, findAddress FindConsumerAddress) (*commonapi.Address, error)
	UpdateAddress(ctx context.Context, updateAddress ModifyConsumerAddress) error
	RemoveAddress(ctx context.Context, removeAddress FindConsumerAddress) error
}
