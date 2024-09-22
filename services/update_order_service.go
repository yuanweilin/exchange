package services

import (
	"encoding/json"
	"exchange/database"
	"exchange/models"
	"fmt"

	"exchange/websocket/controllers"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func FindOrCreateAndUpdateBestBuyOrder(newPrice string) (models.Order, bool, error) {
	newPriceDecimal, _ := decimal.NewFromString(newPrice)

	var order models.Order
	err := database.DB.Where("type = ?", models.OrderTypeBuy).
		Order("price desc").
		First(&order).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			order = models.Order{
				Amount: decimal.NewFromFloat(1.0),
				Price:  newPriceDecimal,
				Type:   models.OrderTypeBuy,
			}
			err = database.DB.Create(&order).Error
			if err != nil {
				return models.Order{}, false, fmt.Errorf("failed to create new buy order: %w", err)
			}
			fmt.Println("Created new buy order")
			return order, true, nil
		} else {
			return models.Order{}, false, fmt.Errorf("failed to find best buy order: %w", err)
		}
	} else {
		if !order.Price.Equal(newPriceDecimal) {
			order.Price = newPriceDecimal
			err = database.DB.Save(&order).Error
			if err != nil {
				return models.Order{}, false, fmt.Errorf("failed to update buy order: %w", err)
			}
			fmt.Println("Updated buy order price")
			return order, true, nil
		} else {
			fmt.Println("Buy order price is the same, no update needed")
			return order, false, nil
		}
	}
}

func FindOrCreateAndUpdateBestSellOrder(newPrice string) (models.Order, bool, error) {
	newPriceDecimal, _ := decimal.NewFromString(newPrice)

	var order models.Order
	err := database.DB.Where("type = ?", models.OrderTypeSell).
		Order("price asc").
		First(&order).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			order = models.Order{
				Amount: decimal.NewFromFloat(1.0),
				Price:  newPriceDecimal,
				Type:   models.OrderTypeSell,
			}
			err = database.DB.Create(&order).Error
			if err != nil {
				return models.Order{}, false, fmt.Errorf("failed to create new sell order: %w", err)
			}
			fmt.Println("Created new sell order")
			return order, true, nil
		} else {
			return models.Order{}, false, fmt.Errorf("failed to find best sell order: %w", err)
		}
	} else {
		if !order.Price.Equal(newPriceDecimal) {
			order.Price = newPriceDecimal
			err = database.DB.Save(&order).Error
			if err != nil {
				return models.Order{}, false, fmt.Errorf("failed to update sell order: %w", err)
			}
			fmt.Println("Updated sell order price")
			return order, true, nil
		} else {
			fmt.Println("Sell order price is the same, no update needed")
			return order, false, nil
		}
	}
}

func UpdateBestOrders(newBuyPrice string, newSellPrice string) error {
	bestBuyOrder, buyOrderChanged, err := FindOrCreateAndUpdateBestBuyOrder(newBuyPrice)
	if err != nil {
		return fmt.Errorf("failed to find or update best buy order: %w", err)
	}

	bestSellOrder, sellOrderChanged, err := FindOrCreateAndUpdateBestSellOrder(newSellPrice)
	if err != nil {
		return fmt.Errorf("failed to find or update best sell order: %w", err)
	}

	if buyOrderChanged {
		broadcastOrderUpdate("buy", bestBuyOrder)
	}

	if sellOrderChanged {
		broadcastOrderUpdate("sell", bestSellOrder)
	}

	fmt.Printf("Best Buy Order: %+v\n", bestBuyOrder)
	fmt.Printf("Best Sell Order: %+v\n", bestSellOrder)

	return nil
}

func broadcastOrderUpdate(orderType string, order models.Order) {
	message := struct {
		Type  string       `json:"type"`
		Order models.Order `json:"order"`
	}{
		Type:  orderType,
		Order: order,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Failed to marshal order update:", err)
		return
	}

	controllers.BroadcastMessage(websocket.TextMessage, messageBytes)
}
