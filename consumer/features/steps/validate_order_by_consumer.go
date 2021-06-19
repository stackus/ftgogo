package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterValidateOrderByConsumerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I validate an order for (?:a|the) consumer with:$`, f.iValidateAnOrderForAConsumerWith)
}

func (f *FeatureState) iValidateAnOrderForAConsumerWith(doc *godog.DocString) error {
	var cmd commands.ValidateOrderByConsumer

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.ConsumerID == "<ConsumerID>" {
		cmd.ConsumerID = f.consumerID
	}

	f.err = f.app.ValidateOrderByConsumer(context.Background(), cmd)

	return nil
}
