package services

import (
	"exchange/database"
	"exchange/models"
	"exchange/websocket/controllers"
	"log"
	"sort"

	"gorm.io/gorm"
)

func fetchOrderBookFromDatabase(db *gorm.DB) models.OrderBook {
	var buyOrders []models.Order
	var sellOrders []models.Order

	if err := db.Where("type = ?", "buy").Find(&buyOrders).Error; err != nil {
		log.Println("Error fetching buy orders:", err)
	}

	if err := db.Where("type = ?", "sell").Find(&sellOrders).Error; err != nil {
		log.Println("Error fetching sell orders:", err)
	}

	sort.Slice(buyOrders, func(i, j int) bool {
		return buyOrders[i].Price.Cmp(buyOrders[j].Price) > 0
	})

	sort.Slice(sellOrders, func(i, j int) bool {
		return sellOrders[i].Price.Cmp(sellOrders[j].Price) < 0
	})

	return models.OrderBook{
		Buy:  buyOrders,
		Sell: sellOrders,
	}
}

// 批次插入新訂單並刪除舊訂單
func replaceOldOrdersWithNew(bestBid string, bestAsk string) {
	db := database.DB
	tx := db.Begin()

	if err := tx.Where("type = ?", "buy").Delete(&models.Order{}).Error; err != nil {
		tx.Rollback()
		return
	}

	if err := tx.Where("type = ?", "sell").Delete(&models.Order{}).Error; err != nil {
		tx.Rollback()
		return
	}

	newBidOrders := GenerateSellOrdersByTraderType(bestBid, models.OrderTypeBuy)
	newAskOrders := GenerateBidOrdersByMarketDepth(bestAsk, models.OrderTypeSell)

	for _, order := range newBidOrders {
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	for _, order := range newAskOrders {
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed to commit transaction:", err)
		return
	}

	orderBook := fetchOrderBookFromDatabase(db)
	controllers.UpdateOrderBook(orderBook)
}
