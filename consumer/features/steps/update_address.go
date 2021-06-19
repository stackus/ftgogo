package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterUpdateAddressSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I update (?:an|the|another) address with:$`, f.iUpdateAnAddressWith)
}

func (f *FeatureState) iUpdateAnAddressWith(doc *godog.DocString) error {
	var cmd commands.UpdateAddress

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.ConsumerID == "<ConsumerID>" {
		cmd.ConsumerID = f.consumerID
	}

	f.err = f.app.UpdateAddress(context.Background(), cmd)

	return nil
}
