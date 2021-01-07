package orderapi

// Service Commands
const OrderServiceCommandChannel = "orderservice"

// Aggregates
const OrderAggregateChannel = "orderservice.Order"

// Sagas
const CreateOrderSagaChannel = "orderservice.CreateOrderSaga"
const CancelOrderSagaChannel = "orderservice.CancelOrderSaga"
const ReviseOrderSagaChannel = "orderservice.ReviseOrderSaga"

func RegisterTypes() {
	registerCommands()
	registerEvents()
	registerReplies()
}
