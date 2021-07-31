package restaurantapi

import (
	restaurantpb "github.com/stackus/ftgogo/serviceapis/restaurantapi/pb"
)

type MenuItem struct {
	ID    string
	Name  string
	Price int
}

func ToMenuItemProto(menuItem MenuItem) *restaurantpb.MenuItem {
	return &restaurantpb.MenuItem{
		ID:    menuItem.ID,
		Name:  menuItem.Name,
		Price: int64(menuItem.Price),
	}
}

func FromMenuItemProto(menuItem *restaurantpb.MenuItem) MenuItem {
	return MenuItem{
		ID:    menuItem.ID,
		Name:  menuItem.Name,
		Price: int(menuItem.Price),
	}
}
