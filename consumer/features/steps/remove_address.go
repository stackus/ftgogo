package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterRemoveAddressSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I remove an address for "([^"]*)" with label "([^"]*)"$`, f.iRemoveAnAddressForWithLabel)
}

func (f *FeatureState) iRemoveAnAddressForWithLabel(consumerName, addressLabel string) error {
	consumerID := f.registeredConsumers[consumerName]

	f.err = f.app.RemoveAddress(context.Background(), commands.RemoveAddress{
		ConsumerID: consumerID,
		AddressID:  addressLabel,
	})

	return nil
}
