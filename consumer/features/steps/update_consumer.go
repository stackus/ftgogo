package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterUpdateConsumerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I update (?:a|the) consumer with:$`, f.iUpdateAConsumerWith)
}

func (f *FeatureState) iUpdateAConsumerWith(doc *godog.DocString) error {
	var cmd commands.UpdateConsumer

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.ConsumerID == "<ConsumerID>" {
		cmd.ConsumerID = f.consumerID
	}

	f.err = f.app.UpdateConsumer(context.Background(), cmd)

	return nil
}
