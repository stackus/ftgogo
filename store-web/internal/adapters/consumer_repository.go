package adapters

import (
	"context"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/store-web/internal/domain"
)

type (
	RegisterConsumer struct {
		Name string
	}

	FindConsumer struct {
		ConsumerID string
	}

	UpdateConsumer struct {
		Consumer domain.Consumer
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
	Find(ctx context.Context, findConsumer FindConsumer) (*domain.Consumer, error)
}
