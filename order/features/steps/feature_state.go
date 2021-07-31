package steps

import (
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/edat/inmem"
	"github.com/stackus/edat/msg"
	"github.com/stackus/errors"

	"github.com/stackus/ftgogo/order/internal/adapters"
	"github.com/stackus/ftgogo/order/internal/application"
	"github.com/stackus/ftgogo/order/internal/domain"
	"github.com/stackus/ftgogo/serviceapis"
)

type FeatureState struct {
	app       application.ServiceApplication
	orderID   string
	orderIDs  map[string]string
	ticketIDs map[string]string
	err       error
}

func NewFeatureState() *FeatureState {
	f := &FeatureState{}
	f.Reset()

	return f
}

var assist = assistdog.NewDefault()

func init() {
	serviceapis.RegisterTypes()
	domain.RegisterTypes()
}

func (f *FeatureState) Reset() {
	f.orderID = ""
	f.orderIDs = make(map[string]string)
	f.ticketIDs = make(map[string]string)
	f.err = nil

	sagaInstanceStore := inmem.NewSagaInstanceStore()
	publisher := msg.NewPublisher(inmem.NewProducer())
	orderRepo := adapters.NewOrderAggregateRepository(inmem.NewEventStore())
	restaurantRepo := adapters.NewRestaurantInmemRepository()
	createOrderSaga := adapters.NewCreateOrderOrchestrationSaga(sagaInstanceStore, publisher)
	cancelOrderSaga := adapters.NewCancelOrderOrchestrationSaga(sagaInstanceStore, publisher)
	reviseOrderSaga := adapters.NewReviseOrderOrchestrationSaga(sagaInstanceStore, publisher)
	ordersPlacedCounter := adapters.NewInmemCounter("orders_placed")
	ordersApprovedCounter := adapters.NewInmemCounter("orders_approved")
	ordersRejectedCounter := adapters.NewInmemCounter("orders_rejected")
	f.app = application.NewServiceApplication(
		orderRepo, restaurantRepo,
		createOrderSaga, cancelOrderSaga, reviseOrderSaga,
		ordersPlacedCounter, ordersApprovedCounter, ordersRejectedCounter,
	)
}

func (f *FeatureState) RegisterCommonSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I expect the (?:request|command|query) to fail$`, f.iExpectTheCommandToFail)
	ctx.Step(`^I expect the (?:request|command|query) to succeed$`, f.iExpectTheCommandToSucceed)

	ctx.Step(`^(?:ensure |expect )?the returned error message is "([^"]*)"$`, f.theReturnedErrorMessageIs)
}

func (f *FeatureState) iExpectTheCommandToFail() error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "expected error to not be nil")
	}
	return nil
}

func (f *FeatureState) iExpectTheCommandToSucceed() error {
	if f.err != nil {
		return errors.Wrapf(f.err, "expected error to be nil: got %w", f.err)
	}

	return nil
}

func (f *FeatureState) theReturnedErrorMessageIs(errorMsg string) error {
	if f.err == nil {
		return errors.Wrap(errors.ErrUnknown, "expected error to not be nil")
	}

	if errorMsg != f.err.Error() {
		return errors.Wrapf(errors.ErrInvalidArgument, "expected: %s: got: %s", errorMsg, f.err.Error())
	}

	return nil
}
