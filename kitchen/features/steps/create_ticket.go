package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
)

func (f *FeatureState) RegisterCreateTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:create|setup|have created) (?:a|the|another) ticket with:$`, f.iCreateATicketWith)
	ctx.Step(`^I (?:cancel|have cancell?ed) creating (?:a|the|another) ticket with:$`, f.iCancelCreateATicketWith)
	ctx.Step(`^I (?:confirm|have confirmed) creating (?:a|the|another) ticket with:$`, f.iConfirmCreateATicketWith)
}

func (f *FeatureState) iCreateATicketWith(doc *godog.DocString) error {
	var cmd commands.CreateTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	f.ticketID, f.err = f.app.CreateTicket(context.Background(), cmd)

	return nil
}

func (f *FeatureState) iCancelCreateATicketWith(doc *godog.DocString) error {
	var cmd commands.CancelCreateTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.TicketID == "<TicketID>" {
		cmd.TicketID = f.ticketID
	}

	f.err = f.app.CancelCreateTicket(context.Background(), cmd)

	return nil
}

func (f *FeatureState) iConfirmCreateATicketWith(doc *godog.DocString) error {
	var cmd commands.ConfirmCreateTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.TicketID == "<TicketID>" {
		cmd.TicketID = f.ticketID
	}

	f.err = f.app.ConfirmCreateTicket(context.Background(), cmd)

	return nil
}
