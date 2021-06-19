package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterDisableAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I disabled? (?:an|the) account with:?$`, f.iDisableAnAccountWith)
}

func (f *FeatureState) iDisableAnAccountWith(doc *godog.DocString) error {
	var cmd commands.DisableAccount

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.err = f.app.DisableAccount(context.Background(), cmd)

	return nil
}
