package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
)

func (f *FeatureState) RegisterCreateDeliverySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:create|setup) (?:a|the|another) delivery with:$`, f.iCreateADeliveryWith)
}

func (f *FeatureState) iCreateADeliveryWith(doc *godog.DocString) error {
	var cmd commands.CreateDelivery

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.err = f.app.CreateDelivery(context.Background(), cmd)

	return nil
}
