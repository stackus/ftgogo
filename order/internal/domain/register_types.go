package domain

func RegisterTypes() {
	registerOrderCommands()
	registerOrderEvents()
	registerOrderSnapshots()
	registerCancelOrderSagaData()
	registerCreateOrderSagaData()
	registerReviseOrderSagaData()
}
