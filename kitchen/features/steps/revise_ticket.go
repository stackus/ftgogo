package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/kitchen/internal/application/commands"
)

func (f *FeatureState) RegisterReviseTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:begin|began|have begun) (?:to revise|revising) (?:a|the|another) ticket with:$`, f.iBeginToReviseATicketWith)
	ctx.Step(`^I (?:undo|have undone) revising (?:a|the|another) ticket with:$`, f.iUndoReviseATicketWith)
	ctx.Step(`^I (?:confirm|have confirmed) revising (?:a|the|another) ticket with:$`, f.iConfirmReviseATicketWith)
}

func (f *FeatureState) iBeginToReviseATicketWith(doc *godog.DocString) error {
	var cmd commands.BeginReviseTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.TicketID == "<TicketID>" {
		cmd.TicketID = f.ticketID
	}

	f.err = f.app.BeginReviseTicket(context.Background(), cmd)

	return nil
}

func (f *FeatureState) iUndoReviseATicketWith(doc *godog.DocString) error {
	var cmd commands.UndoReviseTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.TicketID == "<TicketID>" {
		cmd.TicketID = f.ticketID
	}

	f.err = f.app.UndoReviseTicket(context.Background(), cmd)

	return nil
}

func (f *FeatureState) iConfirmReviseATicketWith(doc *godog.DocString) error {
	var cmd commands.ConfirmReviseTicket

	err := json.Unmarshal([]byte(doc.Content), &cmd)
	if err != nil {
		return err
	}

	if cmd.TicketID == "<TicketID>" {
		cmd.TicketID = f.ticketID
	}

	f.err = f.app.ConfirmReviseTicket(context.Background(), cmd)

	return nil
}
