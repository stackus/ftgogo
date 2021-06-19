package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/queries"
)

func (f *FeatureState) RegisterGetAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:request|get|fetch) (?:an|the) account with:?$`, f.iRequestAnAccountWith)
}

func (f *FeatureState) iRequestAnAccountWith(doc *godog.DocString) error {
	var query queries.GetAccount

	err := json.Unmarshal([]byte(doc.Content), &query)
	if err != nil {
		return err
	}

	f.account, f.err = f.app.GetAccount(context.Background(), query)

	return nil
}
