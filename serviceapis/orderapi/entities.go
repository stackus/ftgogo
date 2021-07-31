package orderapi

import (
	orderpb "github.com/stackus/ftgogo/serviceapis/orderapi/pb"
)

type Address struct {
	Street1 string
	Street2 string
	City    string
	State   string
	Zip     string
}

type OrderDetails struct {
	ConsumerID   string
	RestaurantID string
	LineItems    []LineItem
	OrderTotal   int
}

type LineItem struct {
	MenuItemID string
	Name       string
	Price      int
	Quantity   int
}

func (i LineItem) GetTotal() int {
	return i.Price * i.Quantity
}

// OrderState order status states
type OrderState int

// Valid OrderState values
const (
	UnknownOrderState OrderState = iota
	ApprovalPending
	Approved
	Rejected
	CancelPending
	Cancelled
	RevisionPending
)

func (s OrderState) String() string {
	switch s {
	case ApprovalPending:
		return "ApprovalPending"
	case Approved:
		return "Approved"
	case Rejected:
		return "Rejected"
	case CancelPending:
		return "CancelPending"
	case Cancelled:
		return "Cancelled"
	case RevisionPending:
		return "RevisionPending"
	default:
		return "Unknown"
	}
}

func ToOrderStateProto(orderState OrderState) orderpb.OrderState {
	switch orderState {
	case ApprovalPending:
		return orderpb.OrderState_ApprovalPending
	case Approved:
		return orderpb.OrderState_Approved
	case CancelPending:
		return orderpb.OrderState_CancelPending
	case Cancelled:
		return orderpb.OrderState_Cancelled
	case RevisionPending:
		return orderpb.OrderState_RevisionPending
	case Rejected:
		return orderpb.OrderState_Rejected
	default:
		return orderpb.OrderState_Unknown
	}
}

func FromOrderStateProto(orderState orderpb.OrderState) OrderState {
	switch orderState {
	case orderpb.OrderState_ApprovalPending:
		return ApprovalPending
	case orderpb.OrderState_Approved:
		return Approved
	case orderpb.OrderState_Rejected:
		return Rejected
	case orderpb.OrderState_CancelPending:
		return CancelPending
	case orderpb.OrderState_Cancelled:
		return Cancelled
	case orderpb.OrderState_RevisionPending:
		return RevisionPending
	default:
		return UnknownOrderState
	}
}
