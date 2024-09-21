package services

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adshao/go-binance/v2"
)

func handleDepthEvent(event *binance.WsDepthEvent) {
	if len(event.Bids) > 0 && len(event.Asks) > 0 {
		bestBid := event.Bids[0].Price
		bestAsk := event.Asks[0].Price
		log.Printf("Real-time Best Bid: %s, Best Ask: %s", bestBid, bestAsk)

		fmt.Println("---------------")
		fmt.Println("Best bid:", bestBid)
		fmt.Println("Best ask:", bestAsk)
		fmt.Println("---------------")
	} else {
		log.Println("Bids or Asks are empty, skipping event.")
	}
}

func GetLatestOrderPrices() {
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		handleDepthEvent(event)
	}

	errHandler := func(err error) {
		log.Printf("Error: %v", err)
	}

	doneC, stopC, err := binance.WsDepthServe("ETHBTC", wsDepthHandler, errHandler)
	if err != nil {
		log.Fatal(err)
	}

	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalC:
		log.Println("Received interrupt signal, shutting down gracefully...")
		close(stopC)
	case <-doneC:
		log.Println("WebSocket connection closed.")
	}
}
