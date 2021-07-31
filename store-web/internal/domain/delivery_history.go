package domain

type DeliveryHistory struct {
	ID              string
	AssignedCourier string
	CourierActions  []string
	Status          string
}
