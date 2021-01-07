package restaurantapi

// Service Commands
const RestaurantServiceCommandChannel = "restaurantservice"

// Aggregates
const RestaurantAggregateChannel = "restaurantservice.Restaurant"

func RegisterTypes() {
	registerCommands()
	registerEvents()
	registerReplies()
}
