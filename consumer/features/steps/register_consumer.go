package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterRegisterConsumerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:register|setup) (?:a|the|another) consumer with:$`, f.iRegisterAConsumerWith)
}

func (f *FeatureState) iRegisterAConsumerWith(doc *godog.DocString) error {
	var cmd commands.RegisterConsumer

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.consumerID, f.err = f.app.RegisterConsumer(context.Background(), cmd)

	return nil
}
