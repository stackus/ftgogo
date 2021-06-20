package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/delivery/internal/application/queries"
)

func (f *FeatureState) RegisterGetCourierSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:get|fetch|request) (?:a|the|another) courier with:$`, f.iGetACourierWith)
}

func (f *FeatureState) iGetACourierWith(doc *godog.DocString) error {
	var query queries.GetCourier

	err := json.Unmarshal([]byte(doc.Content), &query)
	if err != nil {
		return err
	}

	if query.CourierID == "<AssignedCourierID>" {
		query.CourierID = f.assignedCourierID
	}

	f.courier, f.err = f.app.GetCourier(context.Background(), query)

	return nil
}
