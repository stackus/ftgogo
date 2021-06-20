package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
)

func (f *FeatureState) RegisterAcceptTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:accept|have accepted) (?:a|the|another) ticket with:$`, f.iAcceptATicketWith)
}

func (f *FeatureState) iAcceptATicketWith(doc *godog.DocString) error {
	var cmd commands.AcceptTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.TicketID == "<TicketID>" {
		cmd.TicketID = f.ticketID
	}

	f.err = f.app.AcceptTicket(context.Background(), cmd)

	return nil
}
