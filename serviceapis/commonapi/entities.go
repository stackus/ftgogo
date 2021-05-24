package commonapi

type Address struct {
	Street1 string
	Street2 string
	City    string
	State   string
	Zip     string
}

type MenuItemQuantities map[string]int
