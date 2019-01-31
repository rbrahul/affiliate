package db

//QueueItem - Each Entry Item for Queue
type QueueItem struct {
	ASIN      string
	ProductID int
}

//QueueItems - List of QueueItem
type QueueItems []QueueItem

//Product - Information of the products
type Product struct {
	ProductID          int
	ASIN               string
	ProductAvailablity map[string]bool
}
