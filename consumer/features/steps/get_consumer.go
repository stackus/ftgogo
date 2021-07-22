package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/application/queries"
)

func (f *FeatureState) RegisterGetConsumerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:get|fetch|request) (?:a|the|some) consumer named "([^"]*)"$`, f.iRequestTheConsumerNamed)

	ctx.Step(`^(?:ensure )?the returned consumer (?:to have|has) the name "([^"]*)"$`, f.theReturnedConsumerHasTheName)
	ctx.Step(`^(?:ensure )?the returned consumer (?:to have|has) (\d+) addresses$`, f.theReturnedConsumerHasAddresses)
	ctx.Step(`^(?:ensure )?the returned consumer (?:to have|has) an address with label "([^"]*)"$`, f.theReturnedConsumerHasAnAddressWithLabel)
}

func (f *FeatureState) iRequestTheConsumerNamed(consumerName string) error {
	consumerID := f.registeredConsumers[consumerName]

	f.consumer, f.err = f.app.GetConsumer(context.Background(), queries.GetConsumer{ConsumerID: consumerID})

	return nil
}

func (f *FeatureState) theReturnedConsumerHasTheName(expected string) error {
	if f.consumer == nil {
		return errors.Wrap(errors.ErrNotFound, "expected consumer to not be nil")
	}

	got := f.consumer.Name()

	if got != expected {
		return errors.Wrapf(errors.ErrInvalidArgument, "name does not match expected: %s: got: %s", expected, got)
	}

	return nil
}

func (f *FeatureState) theReturnedConsumerHasAddresses(expected int) error {
	if f.consumer == nil {
		return errors.Wrap(errors.ErrNotFound, "expected consumer to not be nil")
	}

	got := len(f.consumer.Addresses())

	if got != expected {
		return errors.Wrapf(errors.ErrInvalidArgument, "address count does not match expected: %d: got: %d", expected, got)
	}

	return nil
}

func (f *FeatureState) theReturnedConsumerHasAnAddressWithLabel(addressLabel string) error {
	if f.consumer == nil {
		return errors.Wrap(errors.ErrNotFound, "expected consumer to not be nil")
	}

	if nil == f.consumer.Address(addressLabel) {
		return errors.Wrapf(errors.ErrNotFound, "expected consumer to have an address with label %s", addressLabel)
	}

	return nil
}
