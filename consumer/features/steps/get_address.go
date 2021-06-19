package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/queries"
)

func (f *FeatureState) RegisterGetAddressSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:get|fetch|request) (?:the|an) address with:$`, f.iRequestAnAddressWith)
}

func (f *FeatureState) iRequestAnAddressWith(doc *godog.DocString) error {
	var query queries.GetAddress

	err := json.Unmarshal([]byte(doc.Content), &query)
	if err != nil {
		return err
	}

	if query.ConsumerID == "<ConsumerID>" {
		query.ConsumerID = f.consumerID
	}

	f.address, f.err = f.app.GetAddress(context.Background(), query)

	return nil
}
