package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

func (f *FeatureState) RegisterUpdateAddressSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I update an address for "([^"]*)" with label "([^"]*)"$`, f.iUpdateAnAddressForWithLabel)
}

func (f *FeatureState) iUpdateAnAddressForWithLabel(consumerName, addressLabel string, table *godog.Table) error {
	consumerID := f.registeredConsumers[consumerName]

	address, err := assist.CreateInstance(new(commonapi.Address), table)
	if err != nil {
		return errors.Wrapf(errors.ErrUnknown, "error parsing address table: %w", err)
	}

	f.err = f.app.UpdateAddress(context.Background(), commands.UpdateAddress{
		ConsumerID: consumerID,
		AddressID:  addressLabel,
		Address:    address.(*commonapi.Address),
	})

	return nil
}
