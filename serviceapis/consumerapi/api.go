package consumerapi

// Service Commands
const ConsumerServiceCommandChannel = "consumerservice"

// Aggregates
const ConsumerAggregateChannel = "consumerservice.Consumer"

func RegisterTypes() {
	registerCommands()
	registerEvents()
	registerReplies()
}
