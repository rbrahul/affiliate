package main

import (
	"fmt"
	"sync"

	"github.com/rbrahul/affiliate/db"
	"github.com/rbrahul/affiliate/services"
	"github.com/robfig/cron"
)

var productAvailAbility map[string]bool

func init() {
	productAvailAbility = map[string]bool{
		"us": false,
		"de": false,
		"uk": false,
		"it": false,
		"es": false,
		"fr": false,
		"cn": false,
		"jp": false,
		"in": false,
		"ca": false,
		"br": false,
		"mx": false,
	}
}

func scannASIN(asin string) map[string]bool {
	var wg = sync.WaitGroup{}
	var ProductAvailablity = productAvailAbility
	for code := range services.DomainMapping {
		wg.Add(1)
		go func(countryCode string) {
			ProductAvailablity[countryCode] = services.IsProductAvailable(asin, countryCode)
			wg.Done()
		}(code)
	}

	wg.Wait()
	return ProductAvailablity
}

func processQueue() {
	var Product db.Product
	queueItem, err := db.PickFirstASIN()
	if err != nil {
		return
	}
	Product.ASIN = queueItem.ASIN
	Product.ProductID = queueItem.ProductID
	Product.ProductAvailablity = scannASIN(queueItem.ASIN)
	isProductSaved := db.SaveScannedProduct(Product)
	if isProductSaved {
		db.DeleteQueueItem(Product.ProductID)
	}
	fmt.Println("========", Product.ASIN, "has been synced ========")
	processQueue()
}

func main() {
	c := cron.New()
	c.AddFunc("@every 15m", func() {
		db.SaveASINsIntoQueue()
		processQueue()
		fmt.Println("Successfully Synced")
	})
	c.Start()
	fmt.Scanln()
}
