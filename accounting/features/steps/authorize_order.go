package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/internal/application/commands"
)

func (f *FeatureState) RegisterAuthorizeOrderSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I authorize (?:an|the) order with:?$`, f.iAuthorizeAnOrderWith)
}

func (f *FeatureState) iAuthorizeAnOrderWith(table *godog.Table) error {
	cmd, err := assist.CreateInstance(new(commands.AuthorizeOrder), table)
	if err != nil {
		return err
	}

	f.err = f.app.AuthorizeOrder(context.Background(), *cmd.(*commands.AuthorizeOrder))

	return nil
}
