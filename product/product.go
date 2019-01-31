package product

//Availability - available in countries
type Availability struct {
}

//Product - Information of the products
type Product struct {
	id                 int32
	ASIN               string
	ProductAvailablity map[string]string
	Processed          bool
}
