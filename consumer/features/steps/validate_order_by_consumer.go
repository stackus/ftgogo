package steps

import (
	"context"
	"regexp"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/consumer/internal/application/commands"
)

func (f *FeatureState) RegisterValidateOrderByConsumerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I validate an order for "([^"]*)"$`, f.iValidateAnOrderFor)
}

func (f *FeatureState) iValidateAnOrderFor(consumerName string, table *godog.Table) error {
	type orderData struct {
		OrderID string
		Total   string
	}

	consumerID := f.registeredConsumers[consumerName]

	data, err := assist.CreateInstance(new(orderData), table)
	if err != nil {
		return errors.Wrapf(errors.ErrUnknown, "error parsing order table: %w", err)
	}

	order := data.(*orderData)

	re := regexp.MustCompile(`^\$(\d+)(?:\.(\d+))?`)
	m := re.FindStringSubmatch(order.Total)

	var dollars, cents = 0, 0

	dollars, err = strconv.Atoi(m[1])
	if err != nil {
		return errors.Wrapf(errors.ErrUnknown, "error parsing dollars from total: %s", order.Total)
	}
	cents, err = strconv.Atoi(m[2])
	if err != nil {
		return errors.Wrapf(errors.ErrUnknown, "error parsing cents from total: %s", order.Total)
	}

	cmd := commands.ValidateOrderByConsumer{
		ConsumerID: consumerID,
		OrderID:    order.OrderID,
		OrderTotal: dollars*10 + cents,
	}

	f.err = f.app.ValidateOrderByConsumer(context.Background(), cmd)

	return nil
}
