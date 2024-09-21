package services

import (
	"context"
	"log"

	"github.com/adshao/go-binance/v2"
)

func GetLatestOrderBook(symbol string) (bestBid, bestAsk string, err error) {
	client := binance.NewClient("", "")

	depth, err := client.NewDepthService().Symbol(symbol).Limit(5).Do(context.Background())
	if err != nil {
		log.Println("Error fetching order book:", err)
		return "", "", err
	}

	if len(depth.Bids) > 0 && len(depth.Asks) > 0 {
		bestBid = depth.Bids[0].Price
		bestAsk = depth.Asks[0].Price
		return bestBid, bestAsk, nil
	}

	return "", "", nil
}
