package services

import (
	"math/rand"

	"exchange/models"

	"github.com/shopspring/decimal"
)

func generatePriceNormalDistribution(mean, stdDev float64) float64 {
	r := rand.NormFloat64()
	return mean + r*stdDev
}

func generateTraderOrders(bestAsk decimal.Decimal, traderType string, orderType models.OrderType) models.Order {
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
		Type:   orderType,
	}
}

// 模擬不同類型的交易者行為產生orders
func GenerateSellOrdersByTraderType(bestAsk string, orderType models.OrderType) []models.Order {
	bestAskDecimal, _ := decimal.NewFromString(bestAsk)
	totalAmount := decimal.Zero
	var orders []models.Order

	for totalAmount.LessThan(decimal.NewFromFloat(150.0)) {
		traderType := "institution"
		if rand.Float64() > 0.5 {
			traderType = "retail"
		}

		order := generateTraderOrders(bestAskDecimal, traderType, orderType)
		orders = append(orders, order)

		totalAmount = totalAmount.Add(order.Amount)

		if totalAmount.GreaterThanOrEqual(decimal.NewFromFloat(150.0)) {
			break
		}
	}

	return orders
}

// 基於市場深度的分佈產生orders
func GenerateBidOrdersByMarketDepth(bestBid string, orderType models.OrderType) []models.Order {
	bestBidDecimal, _ := decimal.NewFromString(bestBid)
	totalValue := decimal.Zero
	var orders []models.Order

	meanPrice, _ := bestBidDecimal.Float64()
	stdDev := meanPrice * 0.005

	for totalValue.LessThan(decimal.NewFromFloat(5.0)) {
		price := decimal.NewFromFloat(generatePriceNormalDistribution(meanPrice, stdDev)).Round(4)
		amount := decimal.NewFromFloat(rand.Float64() * 0.5).Round(2)

		if len(orders) == 0 {
			price = bestBidDecimal
		}

		orders = append(orders, models.Order{
			Price:  price,
			Amount: amount,
			Type:   orderType,
		})

		orderValue := amount.Mul(price)
		totalValue = totalValue.Add(orderValue)

		if totalValue.GreaterThanOrEqual(decimal.NewFromFloat(5.0)) {
			break
		}
	}

	return orders
}
