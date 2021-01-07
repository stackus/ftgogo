package orderapi

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
