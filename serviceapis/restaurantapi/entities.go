package restaurantapi

type Address struct {
	Street1 string
	Street2 string
	City    string
	State   string
	Zip     string
}

type MenuItem struct {
	ID    string
	Name  string
	Price int
}
