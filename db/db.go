package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	session, err := sql.Open("mysql", "root:mysql@/wordpress-cloned?charset=utf8&parseTime=True")
	//session, err := sql.Open("mysql", "root:mysql@/go_db?charset=utf8&parseTime=True")
	checkErr(err)
	db = session
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//SaveASINsIntoQueue - Saving to the queue
func SaveASINsIntoQueue() {
	_, errs := db.Query("INSERT INTO `product-sync-queue`(`product_id`,`ASIN`) select wp.ID, wpm.meta_value from wp_posts wp INNER JOIN `wp_postmeta` wpm ON wp.ID = wpm.`post_id` AND wp.`post_type` = 'product' AND wpm.`meta_key`='_amzASIN' AND wp.ID NOT IN (SELECT `product_id` FROM `scanned-products`)")
	checkErr(errs)
}

//PickFirstASIN - Pick first asin
func PickFirstASIN() (QueueItem, error) {
	var QueueItem QueueItem
	errs := db.QueryRow("SELECT `product_id`, `asin` FROM `product-sync-queue` LIMIT 1").Scan(&QueueItem.ProductID, &QueueItem.ASIN)
	if errs != nil {
		fmt.Println(errs)
		return QueueItem, errs
	}
	return QueueItem, nil
}

//SaveScannedProduct - Save Scanned Product
func SaveScannedProduct(product Product) bool {
	availability, _ := json.Marshal(product.ProductAvailablity)
	_, errs := db.Query("INSERT INTO `scanned-products`(`product_id`,`asin`,`available_countries`) VALUES(?,?,?)", product.ProductID, product.ASIN, availability)
	if errs != nil {
		fmt.Println(errs)
		return false
	}
	return true
}

//DeleteQueueItem - DeleteQueueItem  FROM QUE ITEM BY PRODUCT ID
func DeleteQueueItem(productID int) bool {
	_, errs := db.Query("DELETE FROM `product-sync-queue` WHERE `product_id`=?", productID)
	if errs != nil {
		fmt.Println(errs)
		return false
	}
	return true
}
