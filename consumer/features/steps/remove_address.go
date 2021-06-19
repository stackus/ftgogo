package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterRemoveAddressSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I remove (?:an|the|another) address with:$`, f.iRemoveAnAddressWith)
}

func (f *FeatureState) iRemoveAnAddressWith(doc *godog.DocString) error {
	var cmd commands.RemoveAddress

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.ConsumerID == "<ConsumerID>" {
		cmd.ConsumerID = f.consumerID
	}

	f.err = f.app.RemoveAddress(context.Background(), cmd)

	return nil
}
