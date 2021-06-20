package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/delivery/internal/application/queries"
)

func (f *FeatureState) RegisterGetDeliverySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:get|fetch|request) (?:a|the|another) delivery with:$`, f.iGetADeliveryWith)
}

func (f *FeatureState) iGetADeliveryWith(doc *godog.DocString) error {
	var query queries.GetDelivery

	err := json.Unmarshal([]byte(doc.Content), &query)
	if err != nil {
		return err
	}

	f.delivery, f.err = f.app.GetDelivery(context.Background(), query)
	if f.delivery != nil {
		f.assignedCourierID = f.delivery.AssignedCourierID
	}

	return nil
}
