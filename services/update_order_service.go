package services

import (
	"exchange/database"
	"exchange/models"
	"exchange/websocket"
	"log"
	"sort"
)

func FetchOrderBookFromDatabase() models.OrderBook {
	db := database.DB
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

	buyOrders = mergeOrdersWithSamePrice(buyOrders)
	sellOrders = mergeOrdersWithSamePrice(sellOrders)

	return models.OrderBook{
		Buy:  buyOrders,
		Sell: sellOrders,
	}
}

func mergeOrdersWithSamePrice(orders []models.Order) []models.Order {
	if len(orders) == 0 {
		return orders
	}

	mergedOrders := []models.Order{}
	currentOrder := orders[0]

	for i := 1; i < len(orders); i++ {
		if orders[i].Price.Equal(currentOrder.Price) {
			currentOrder.Amount = currentOrder.Amount.Add(orders[i].Amount)
		} else {
			mergedOrders = append(mergedOrders, currentOrder)
			currentOrder = orders[i]
		}
	}
	mergedOrders = append(mergedOrders, currentOrder)

	return mergedOrders
}

// 批次插入新訂單並刪除舊訂單
func replaceOldOrdersWithNew(bestBid string, bestAsk string) {
	db := database.DB
	tx := db.Begin()

	newBidOrders := GenerateBidOrdersByMarketDepth(bestBid)
	newAskOrders := GenerateSellOrdersByTraderType(bestAsk)

	if err := tx.Where("type = ?", "buy").Delete(&models.Order{}).Error; err != nil {
		tx.Rollback()
		return
	}

	if err := tx.Where("type = ?", "sell").Delete(&models.Order{}).Error; err != nil {
		tx.Rollback()
		return
	}

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

	orderBook := FetchOrderBookFromDatabase()
	websocket.UpdateOrderBook(orderBook)
}
