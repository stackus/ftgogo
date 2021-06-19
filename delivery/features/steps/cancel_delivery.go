package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
)

func (f *FeatureState) RegisterCancelDeliverySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I cancel (?:a|the) delivery with:$`, f.iCancelADeliveryWith)
}

func (f *FeatureState) iCancelADeliveryWith(doc *godog.DocString) error {
	var cmd commands.CancelDelivery

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.err = f.app.CancelDelivery(context.Background(), cmd)

	return nil
}
