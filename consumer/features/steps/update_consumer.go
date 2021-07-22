package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterUpdateConsumerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I change "([^"]*)" name to "([^"]*)"$`, f.iChangeNameTo)
}

func (f *FeatureState) iChangeNameTo(currentName, newName string) error {
	consumerID := f.registeredConsumers[currentName]

	f.err = f.app.UpdateConsumer(context.Background(), commands.UpdateConsumer{
		ConsumerID: consumerID,
		Name:       newName,
	})

	if f.err == nil {
		delete(f.registeredConsumers, currentName)
		f.registeredConsumers[newName] = consumerID
	}

	return nil
}
