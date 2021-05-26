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
	RegisterConsumer struct {
		Name string
	}

	FindConsumer struct {
		ConsumerID string
	}

	UpdateConsumer struct {
		Consumer Consumer
	}

	AddConsumerAddress struct {
		ConsumerID string
		AddressID  string
		Address    *commonapi.Address
	}

	FindConsumerAddress struct {
		ConsumerID string
		AddressID  string
	}

	UpdateConsumerAddress AddConsumerAddress
	RemoveConsumerAddress FindConsumerAddress
)

type ConsumerRepository interface {
	Register(ctx context.Context, registerConsumer RegisterConsumer) (string, error)
	Find(ctx context.Context, findConsumer FindConsumer) (*Consumer, error)
	Update(ctx context.Context, updateConsumer UpdateConsumer) error
	AddAddress(ctx context.Context, addAddress AddConsumerAddress) error
	FindAddress(ctx context.Context, findAddress FindConsumerAddress) (*commonapi.Address, error)
	UpdateAddress(ctx context.Context, updateAddress UpdateConsumerAddress) error
	RemoveAddress(ctx context.Context, removeAddress RemoveConsumerAddress) error
}
