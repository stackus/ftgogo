package kitchenapi

// Service Commands
const KitchenServiceCommandChannel = "kitchenservice"

// Aggregates
const TicketAggregateChannel = "kitchenservice.Ticket"

func RegisterTypes() {
	registerCommands()
	registerEvents()
	registerReplies()
}
