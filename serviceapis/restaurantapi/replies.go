package restaurantapi

import (
	"github.com/stackus/edat/core"
)

func registerReplies() {
	core.RegisterReplies(MenuItemsValidated{})
}

type MenuItemsValidated struct {
	RestaurantID string
	MenuItems    map[string]struct {
		Name  string
		Price int
	}
}

func (MenuItemsValidated) ReplyName() string { return "restaurantapi.MenuItemsValidated" }
