package steps

import (
	"context"
	"reflect"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

func (f *FeatureState) RegisterGetAddressSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I request the "([^"]*)" address for "([^"]*)"$`, f.iRequestTheAddressFor)

	ctx.Step(`^(?:ensure )?the returned address (?:to match|matches)$`, f.theReturnedAddressToMatch)
}

func (f *FeatureState) iRequestTheAddressFor(addressLabel, consumerName string) error {
	consumerID := f.registeredConsumers[consumerName]

	f.address, f.err = f.app.GetAddress(context.Background(), queries.GetAddress{
		ConsumerID: consumerID,
		AddressID:  addressLabel,
	})

	return nil
}

func (f *FeatureState) theReturnedAddressToMatch(table *godog.Table) error {
	expected, err := assist.CreateInstance(new(commonapi.Address), table)
	if err != nil {
		return errors.Wrapf(errors.ErrUnknown, "error parsing address table: %w", err)
	}

	if f.address == nil {
		return errors.Wrap(errors.ErrNotFound, "expected address to not be nil")
	}

	if !reflect.DeepEqual(expected, f.address) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, f.address)
	}

	return nil
}
