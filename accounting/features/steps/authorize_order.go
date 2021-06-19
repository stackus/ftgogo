package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterAuthorizeOrderSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I authorize (?:an|the) order with:?$`, f.iAuthorizeAnOrderWith)
}

func (f *FeatureState) iAuthorizeAnOrderWith(doc *godog.DocString) error {
	var cmd commands.AuthorizeOrder

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.err = f.app.AuthorizeOrder(context.Background(), cmd)

	return nil
}
