package steps

import (
	"context"
	"encoding/json"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/kitchen/internal/application/queries"
)

func (f *FeatureState) RegisterGetTicketSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:get|fetch|request) (?:a|the|another) ticket with:$`, f.iGetATicketWith)
}

func (f *FeatureState) iGetATicketWith(doc *godog.DocString) error {
	var query queries.GetTicket

	err := json.Unmarshal([]byte(doc.Content), &query)
	if err != nil {
		return err
	}

	if query.TicketID == "<TicketID>" {
		query.TicketID = f.ticketID
	}

	f.ticket, f.err = f.app.GetTicket(context.Background(), query)

	return nil
}
