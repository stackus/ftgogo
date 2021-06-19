package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterCreateAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I create an account with:?$`, f.iCreateAnAccountWith)
}

func (f *FeatureState) iCreateAnAccountWith(doc *godog.DocString) error {
	var cmd commands.CreateAccount

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.err = f.app.CreateAccount(context.Background(), cmd)

	return nil
}
