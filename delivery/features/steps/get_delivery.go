package steps

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/cucumber/godog"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/delivery/internal/application/queries"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

func (f *FeatureState) RegisterGetDeliverySteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I (?:get|fetch|request) the delivery information for order "([^"]*)"$`, f.iRequestDeliveryInformationForOrder)

	ctx.Step(`^(?:ensure |expect )?the returned delivery status is "([^"]*)"$`, f.theReturnedDeliveryStatusIs)
	ctx.Step(`^(?:ensure |expect )?the returned delivery matches:$`, f.theReturnedDeliveryMatches)
	ctx.Step(`^(?:ensure |expect )?the returned delivery is not assigned to:$`, f.theReturnedDeliveryIsNotAssignedTo)
}

func (f *FeatureState) iRequestDeliveryInformationForOrder(orderID string) error {
	f.delivery, f.err = f.app.GetDelivery(context.Background(), queries.GetDelivery{OrderID: orderID})

	if f.delivery != nil {
		f.assignedCourierID = f.delivery.AssignedCourierID
	}

	return nil
}

func (f *FeatureState) theReturnedDeliveryStatusIs(status string) error {
	if f.delivery == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected delivery to not be nil")
	}

	if f.delivery.Status.String() != status {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", status, f.delivery.Status.String())
	}

	return nil
}

func (f *FeatureState) theReturnedDeliveryMatches(doc *godog.DocString) error {
	var err error
	var expected *commonapi.Address

	if err = json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return errors.Wrap(errors.ErrUnprocessableEntity, err.Error())
	}

	if f.delivery == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected delivery to not be nil")
	}

	if !reflect.DeepEqual(expected, f.delivery) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, f.delivery)
	}

	return nil
}

func (f *FeatureState) theReturnedDeliveryIsNotAssignedTo(doc *godog.DocString) error {
	if f.delivery == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected delivery to not be nil")
	}

	if f.delivery.AssignedCourierID == doc.Content {
		return errors.Wrapf(errors.ErrInvalidArgument, "does match expected: %v: got: %v", doc.Content, f.delivery.Status.String())
	}

	return nil
}
