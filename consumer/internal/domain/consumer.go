package domain

import (
	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/serviceapis/commonapi"
	"github.com/stackus/ftgogo/serviceapis/consumerapi"
)

var (
	ErrConsumerUnhandledCommand  = errors.Wrap(errors.ErrInternal, "unhandled command in consumer aggregate")
	ErrConsumerUnhandledEvent    = errors.Wrap(errors.ErrInternal, "unhandled event in consumer aggregate")
	ErrConsumerUnhandledSnapshot = errors.Wrap(errors.ErrInternal, "unhandled snapshot in consumer aggregate")
	ErrConsumerNotFound          = errors.Wrap(errors.ErrNotFound, "consumer not found")
	ErrConsumerNameMissing       = errors.Wrap(errors.ErrInvalidArgument, "cannot register a consumer without a name")
	ErrOrderNotValidated         = errors.Wrap(errors.ErrBadRequest, "order not validated for consumer")
	ErrAddressAlreadyExists      = errors.Wrap(errors.ErrConflict, "address with that identifier already exists")
	ErrAddressDoesNotExist       = errors.Wrap(errors.ErrNotFound, "address with that identifier does not exist")
)

type Consumer struct {
	es.AggregateBase
	name      string
	addresses map[string]*commonapi.Address
}

var _ es.Aggregate = (*Consumer)(nil)

func NewConsumer() es.Aggregate {
	return &Consumer{
		addresses: map[string]*commonapi.Address{},
	}
}

func (Consumer) EntityName() string {
	return "consumerservice.Consumer"
}

func (c *Consumer) Name() string {
	return c.name
}

func (c *Consumer) Addresses() map[string]*commonapi.Address {
	return c.addresses
}

func (c *Consumer) Address(addressID string) *commonapi.Address {
	return c.addresses[addressID]
}

// ValidateOrderByConsumer domain method
func (c *Consumer) ValidateOrderByConsumer(orderTotal int) error {
	// ftgo: implement some business logic
	return nil
}

// ProcessCommand aggregate method
func (c *Consumer) ProcessCommand(command core.Command) error {
	switch cmd := command.(type) {
	case *RegisterConsumer:
		if len([]rune(cmd.Name)) == 0 {
			return ErrConsumerNameMissing
		}

		c.AddEvent(&consumerapi.ConsumerRegistered{
			Name: cmd.Name,
		})

	case *UpdateConsumer:
		c.AddEvent(&consumerapi.ConsumerUpdated{
			Name: cmd.Name,
		})

	case *AddAddress:
		if _, exists := c.addresses[cmd.AddressID]; exists {
			return ErrAddressAlreadyExists
		}

		c.AddEvent(&consumerapi.AddressAdded{
			AddressID: cmd.AddressID,
			Address:   cmd.Address,
		})

	case *UpdateAddress:
		if _, exists := c.addresses[cmd.AddressID]; !exists {
			return ErrAddressDoesNotExist
		}

		c.AddEvent(&consumerapi.AddressUpdated{
			AddressID: cmd.AddressID,
			Address:   cmd.Address,
		})

	case *RemoveAddress:
		if _, exists := c.addresses[cmd.AddressID]; !exists {
			return ErrAddressDoesNotExist
		}

		c.AddEvent(&consumerapi.AddressRemoved{
			AddressID: cmd.AddressID,
		})

	default:
		return errors.Wrapf(ErrConsumerUnhandledCommand, "unhandled command `%s`", command.CommandName())
	}

	return nil
}

// ApplyEvent aggregate method
func (c *Consumer) ApplyEvent(event core.Event) error {
	switch evt := event.(type) {
	case *consumerapi.ConsumerRegistered:
		c.name = evt.Name

	case *consumerapi.ConsumerUpdated:
		c.name = evt.Name

	case *consumerapi.AddressAdded:
		c.addresses[evt.AddressID] = evt.Address

	case *consumerapi.AddressUpdated:
		c.addresses[evt.AddressID] = evt.Address

	case *consumerapi.AddressRemoved:
		delete(c.addresses, evt.AddressID)

	default:
		return errors.Wrapf(ErrConsumerUnhandledEvent, "unhandled event `%s`", event.EventName())
	}

	return nil
}

// ApplySnapshot aggregate method
func (c *Consumer) ApplySnapshot(snapshot core.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ConsumerSnapshot:
		c.name = ss.Name
		c.addresses = ss.Addresses

	default:
		return errors.Wrapf(ErrConsumerUnhandledSnapshot, "unhandled snapshot `%s`", snapshot.SnapshotName())
	}

	return nil
}

// ToSnapshot aggregate method
func (c *Consumer) ToSnapshot() (core.Snapshot, error) {
	return &ConsumerSnapshot{
		Name:      c.name,
		Addresses: c.addresses,
	}, nil
}
