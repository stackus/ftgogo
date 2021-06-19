package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterEnableAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I enabled? (?:an|the) account with:?$`, f.iEnableAnAccountWith)
}

func (f *FeatureState) iEnableAnAccountWith(doc *godog.DocString) error {
	var cmd commands.EnableAccount

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.err = f.app.EnableAccount(context.Background(), cmd)

	return nil
}
