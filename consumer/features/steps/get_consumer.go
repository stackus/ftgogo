package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/consumer/internal/application/queries"
)

func (f *FeatureState) RegisterGetConsumerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:get|fetch|request) (?:the|a) consumer with:$`, f.iRequestAnConsumerWith)
}

func (f *FeatureState) iRequestAnConsumerWith(doc *godog.DocString) error {
	var query queries.GetConsumer

	err := json.Unmarshal([]byte(doc.Content), &query)
	if err != nil {
		return err
	}

	if query.ConsumerID == "<ConsumerID>" {
		query.ConsumerID = f.consumerID
	}

	f.consumer, f.err = f.app.GetConsumer(context.Background(), query)

	return nil
}
