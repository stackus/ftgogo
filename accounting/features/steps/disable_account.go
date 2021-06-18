package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterDisableAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I disabled? (?:an|the) account with:?$`, f.iDisableAnAccountWith)
}

func (f *FeatureState) iDisableAnAccountWith(table *godog.Table) error {
	cmd, err := assist.CreateInstance(new(commands.DisableAccount), table)
	if err != nil {
		return err
	}

	f.err = f.app.DisableAccount(context.Background(), *cmd.(*commands.DisableAccount))

	return nil
}
