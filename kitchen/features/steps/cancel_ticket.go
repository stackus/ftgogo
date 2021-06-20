package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
)

func (f *FeatureState) RegisterCancelTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:begin|began|have begun) (?:to cancel|cancelling) (?:a|the|another) ticket with:$`, f.iBeginToCancelATicketWith)
	ctx.Step(`^I (?:undo|have undone) cancell?ing (?:a|the|another) ticket with:$`, f.iUndoCancelATicketWith)
	ctx.Step(`^I (?:confirm|have confirmed) cancell?ing (?:a|the|another) ticket with:$`, f.iConfirmCancelATicketWith)
}

func (f *FeatureState) iBeginToCancelATicketWith(doc *godog.DocString) error {
	var cmd commands.BeginCancelTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.TicketID == "<TicketID>" {
		cmd.TicketID = f.ticketID
	}

	f.err = f.app.BeginCancelTicket(context.Background(), cmd)

	return nil
}

func (f *FeatureState) iUndoCancelATicketWith(doc *godog.DocString) error {
	var cmd commands.UndoCancelTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.TicketID == "<TicketID>" {
		cmd.TicketID = f.ticketID
	}

	f.err = f.app.UndoCancelTicket(context.Background(), cmd)

	return nil
}

func (f *FeatureState) iConfirmCancelATicketWith(doc *godog.DocString) error {
	var cmd commands.ConfirmCancelTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.TicketID == "<TicketID>" {
		cmd.TicketID = f.ticketID
	}

	f.err = f.app.ConfirmCancelTicket(context.Background(), cmd)

	return nil
}
