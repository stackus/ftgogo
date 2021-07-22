package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterRegisterConsumerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I register (?:a|another) consumer named "([^"]*)"$`, f.iRegisterAConsumerNamed)
}

func (f *FeatureState) iRegisterAConsumerNamed(consumerName string) error {
	cmd := commands.RegisterConsumer{
		Name: consumerName,
	}

	f.consumerID, f.err = f.app.RegisterConsumer(context.Background(), cmd)
	f.registeredConsumers[consumerName] = f.consumerID

	return nil
}
