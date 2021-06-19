package steps

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cucumber/godog"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/delivery/internal/adapters"
	"github.com/stackus/ftgogo/delivery/internal/application"
	"github.com/stackus/ftgogo/delivery/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
	"github.com/stackus/ftgogo/serviceapis/commonapi"
)

type ConsumerJson struct {
	ID        string
	Name      string
	Addresses map[string]*commonapi.Address
}

type FeatureState struct {
	app               application.ServiceApplication
	courier           *domain.Courier
	delivery          *domain.Delivery
	assignedCourierID string
	err               error
}

func NewFeatureState() *FeatureState {
	f := &FeatureState{}
	f.Reset()

	return f
}

func init() {
	serviceapis.RegisterTypes()
}

func (f *FeatureState) Reset() {
	f.courier = nil
	f.delivery = nil
	f.assignedCourierID = ""
	f.err = nil

	courierRepo := adapters.NewCourierInmemRepository()
	deliveryRepo := adapters.NewDeliveryInmemRepository()
	restaurantRepo := adapters.NewRestaurantInmemRepository()
	f.app = application.NewServiceApplication(courierRepo, deliveryRepo, restaurantRepo)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure )?the returned courier matches:$`, f.theReturnedCourierMatches)
	ctx.Step(`^(?:ensure )?the returned delivery matches:$`, f.theReturnedDeliveryMatches)
	ctx.Step(`^(?:ensure )?the returned error message is:$`, f.theReturnedErrorMessageIs)
	ctx.Step(`^(?:ensure )?the returned delivery status is:$`, f.theReturnedDeliveryStatusIs)
	ctx.Step(`^(?:ensure )?the returned delivery is not assigned to:$`, f.theReturnedDeliveryIsNotAssignedTo)
	ctx.Step(`^(?:ensure )?the returned courier is (not(?: ))?available$`, f.theReturnedCourierIsAvailable)
}

func (f *FeatureState) iExpectTheCommandToFail() error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "Expected error to not be nil")
	}
	return nil
}

func (f *FeatureState) iExpectTheCommandToSucceed() error {
	if f.err != nil {
		return errors.Wrap(f.err, "Expected error to be nil")
	}
	return nil
}

func (f *FeatureState) theReturnedCourierMatches(doc *godog.DocString) error {
	var err error
	var expected *domain.Courier

	if err = json.Unmarshal([]byte(doc.Content), &expected); err != nil {
		return errors.Wrap(errors.ErrUnprocessableEntity, err.Error())
	}

	if f.courier == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected courier to not be nil")
	}

	if expected.CourierID == "<AssignedCourierID>" {
		expected.CourierID = f.assignedCourierID
	}

	if !reflect.DeepEqual(expected, f.courier) {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", expected, f.courier)
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

func (f *FeatureState) theReturnedErrorMessageIs(doc *godog.DocString) error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "Expected error to not be nil")
	}

	if doc.Content != f.err.Error() {
		return errors.Wrapf(errors.ErrInvalidArgument, "expected: %s: got: %s", doc.Content, f.err.Error())
	}

	return nil
}

func (f *FeatureState) theReturnedDeliveryStatusIs(doc *godog.DocString) error {
	if f.delivery == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected delivery to not be nil")
	}

	if f.delivery.Status.String() != doc.Content {
		return errors.Wrapf(errors.ErrInvalidArgument, "does not match expected: %v: got: %v", doc.Content, f.delivery.Status.String())
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

func (f *FeatureState) theReturnedCourierIsAvailable(neg string) error {
	if f.courier == nil {
		return errors.Wrap(errors.ErrNotFound, "Expected courier to not be nil")
	}

	fmt.Println("CHECKING", neg, f.courier)

	if neg == "not" && f.courier.Available {
		return errors.Wrap(errors.ErrInvalidArgument, "expected courier to not be available")
	} else if neg == "" && !f.courier.Available {
		return errors.Wrap(errors.ErrInvalidArgument, "expected courier to be available")
	}

	return nil
}
