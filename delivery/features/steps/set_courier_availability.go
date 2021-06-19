package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/delivery/internal/application/commands"
)

func (f *FeatureState) RegisterSetCourierAvailabilitySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:create|setup) (?:a|the|another) courier with:$`, f.iSetACourierAvailabilityWith)
	ctx.Step(`^I set (?:a|the) couriers? availability with:$`, f.iSetACourierAvailabilityWith)
}

func (f *FeatureState) iSetACourierAvailabilityWith(doc *godog.DocString) error {
	var cmd commands.SetCourierAvailability

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.err = f.app.SetCourierAvailability(context.Background(), cmd)

	return nil
}
