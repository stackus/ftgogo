package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterAddAddressSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I add (?:an|the|another) address with:$`, f.iAddAnAddressWith)
}

func (f *FeatureState) iAddAnAddressWith(doc *godog.DocString) error {
	var cmd commands.AddAddress

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.ConsumerID == "<ConsumerID>" {
		cmd.ConsumerID = f.consumerID
	}

	f.err = f.app.AddAddress(context.Background(), cmd)

	return nil
}
