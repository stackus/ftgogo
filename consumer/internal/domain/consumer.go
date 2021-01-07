package domain

import (
	"fmt"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/es"
	"serviceapis/consumerapi"
	"shared-go/errs"
)

var ErrConsumerUnhandledCommand = errs.NewError("unhandled command in domain aggregate", errs.ErrServerError)
var ErrConsumerUnhandledEvent = errs.NewError("unhandled event in domain aggregate", errs.ErrServerError)
var ErrConsumerUnhandledSnapshot = errs.NewError("unhandled snapshot in domain aggregate", errs.ErrServerError)

var ErrConsumerNotFound = errs.NewError("domain not found", errs.ErrNotFound)
var ErrOrderNotValidated = errs.NewError("order not validated for domain", errs.ErrBadRequest)

type Consumer struct {
	es.AggregateBase
	name string
}

var _ es.Aggregate = (*Consumer)(nil)

func NewConsumer() es.Aggregate {
	return &Consumer{}
}

func (Consumer) EntityName() string {
	return "consumerservice.Consumer"
}

func (c *Consumer) Name() string {
	return c.name
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
		c.AddEvent(&consumerapi.ConsumerRegistered{
			Name: cmd.Name,
		})

	default:
		return errs.NewError(fmt.Sprintf("unhandled command `%T`", command), ErrConsumerUnhandledCommand)
	}

	return nil
}

// ApplyEvent aggregate method
func (c *Consumer) ApplyEvent(event core.Event) error {
	switch evt := event.(type) {
	case *consumerapi.ConsumerRegistered:
		c.name = evt.Name

	default:
		return errs.NewError(fmt.Sprintf("unhandled event `%T`", event), ErrConsumerUnhandledEvent)
	}

	return nil
}

// ApplySnapshot aggregate method
func (c *Consumer) ApplySnapshot(snapshot core.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ConsumerSnapshot:
		c.name = ss.Name

	default:
		return errs.NewError(fmt.Sprintf("unhandled snapshot `%T`", snapshot), ErrConsumerUnhandledSnapshot)
	}

	return nil
}

// ToSnapshot aggregate method
func (c *Consumer) ToSnapshot() (core.Snapshot, error) {
	return &ConsumerSnapshot{
		Name: c.name,
	}, nil
}
