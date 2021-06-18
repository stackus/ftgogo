package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterCreateAccountSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I create an account with:?$`, f.iCreateAnAccountWith)
}

func (f *FeatureState) iCreateAnAccountWith(table *godog.Table) error {
	cmd, err := assist.CreateInstance(new(commands.CreateAccount), table)
	if err != nil {
		return err
	}

	f.err = f.app.CreateAccount(context.Background(), *cmd.(*commands.CreateAccount))

	return nil
}
