package services

import (
	"log"
	"math/rand"

	"exchange/models"

	"github.com/shopspring/decimal"
)

func generatePriceNormalDistribution(mean, stdDev float64) float64 {
	r := rand.NormFloat64()
	return mean + r*stdDev
}

func generateTraderOrders(bestAsk decimal.Decimal, traderType string) models.Order {
	amount := decimal.NewFromFloat(0)
	price := bestAsk

	if traderType == "institution" {
		amount = decimal.NewFromFloat(rand.Float64()*10 + 5).Round(2)
		price = bestAsk.Mul(decimal.NewFromFloat(1 + rand.Float64()*0.001)).Round(4)
	} else {
		amount = decimal.NewFromFloat(rand.Float64()).Round(2)
		price = bestAsk.Mul(decimal.NewFromFloat(1 + rand.Float64()*0.005)).Round(4)
	}

	return models.Order{
		Price:  price,
		Amount: amount,
		Type:   "sell",
	}
}

// 模擬不同類型的交易者行為產生賣單，保留傳入的最佳賣單價格
func GenerateSellOrdersByTraderType(bestAsk string) []models.Order {
	bestAskDecimal, _ := decimal.NewFromString(bestAsk)
	totalAmount := decimal.Zero
	var orders []models.Order

	orders = append(orders, models.Order{
		Price:  bestAskDecimal,
		Amount: decimal.NewFromFloat(rand.Float64() * 10).Round(2),
		Type:   "sell",
	})
	totalAmount = totalAmount.Add(orders[0].Amount)

	for totalAmount.LessThan(decimal.NewFromFloat(150.0)) {
		traderType := "institution"
		if rand.Float64() > 0.5 {
			traderType = "retail"
		}

		order := generateTraderOrders(bestAskDecimal, traderType)
		if totalAmount.Add(order.Amount).GreaterThan(decimal.NewFromFloat(150.0)) {
			remainingAmount := decimal.NewFromFloat(150.0).Sub(totalAmount)
			order.Amount = remainingAmount.Round(2)
			orders = append(orders, order)
			totalAmount = totalAmount.Add(order.Amount)
			break
		} else {
			orders = append(orders, order)
			totalAmount = totalAmount.Add(order.Amount)
		}
	}

	return orders
}

// 基於市場深度的分佈產生買單，保留傳入的最佳買單價格
func GenerateBidOrdersByMarketDepth(bestBid string) []models.Order {
	bestBidDecimal, _ := decimal.NewFromString(bestBid)
	totalValue := decimal.Zero
	var orders []models.Order

	currentPrice := bestBidDecimal

	amount := decimal.NewFromFloat(rand.Float64()*0.5 + 0.1).Round(2)
	orders = append(orders, models.Order{
		Price:  currentPrice,
		Amount: amount,
		Type:   "buy",
	})
	totalValue = totalValue.Add(amount.Mul(currentPrice))

	meanPrice, _ := bestBidDecimal.Float64()
	stdDev := meanPrice * 0.005

	for totalValue.LessThan(decimal.NewFromFloat(5.0)) {
		price := decimal.NewFromFloat(generatePriceNormalDistribution(meanPrice, stdDev)).Round(4)

		if price.GreaterThanOrEqual(currentPrice) {
			price = currentPrice.Sub(decimal.NewFromFloat(0.0001))
		}
		if price.LessThanOrEqual(decimal.NewFromFloat(0)) {
			log.Println("Price is too low or negative, stopping order generation.")
			break
		}

		amount := decimal.NewFromFloat(rand.Float64()*1 + 0.2).Round(2)
		if amount.LessThanOrEqual(decimal.NewFromFloat(0)) {
			amount = decimal.NewFromFloat(0.1)
		}

		orderValue := amount.Mul(price)
		if totalValue.Add(orderValue).GreaterThan(decimal.NewFromFloat(5.0)) {
			remainingValue := decimal.NewFromFloat(5.0).Sub(totalValue)
			adjustedAmount := remainingValue.Div(price).Round(2)
			if adjustedAmount.GreaterThan(decimal.NewFromFloat(0)) {
				orders = append(orders, models.Order{
					Price:  price,
					Amount: adjustedAmount,
					Type:   "buy",
				})
				totalValue = totalValue.Add(adjustedAmount.Mul(price))
			}
			break
		}

		orders = append(orders, models.Order{
			Price:  price,
			Amount: amount,
			Type:   "buy",
		})

		totalValue = totalValue.Add(orderValue)
		currentPrice = price
		if totalValue.GreaterThanOrEqual(decimal.NewFromFloat(4.5)) {
			break
		}
	}

	return orders
}
