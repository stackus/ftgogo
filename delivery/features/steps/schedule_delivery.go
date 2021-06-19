package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
)

func (f *FeatureState) RegisterScheduleDeliverySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I schedule (?:a|the) delivery with:$`, f.iScheduleADeliveryWith)
}

func (f *FeatureState) iScheduleADeliveryWith(doc *godog.DocString) error {
	var cmd commands.ScheduleDelivery

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.err = f.app.ScheduleDelivery(context.Background(), cmd)

	return nil
}
