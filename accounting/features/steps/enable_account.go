package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterEnableAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I enabled? (?:an|the) account with:?$`, f.iEnableAnAccountWith)
}

func (f *FeatureState) iEnableAnAccountWith(table *godog.Table) error {
	cmd, err := assist.CreateInstance(new(commands.EnableAccount), table)
	if err != nil {
		return err
	}

	f.err = f.app.EnableAccount(context.Background(), *cmd.(*commands.EnableAccount))

	return nil
}
