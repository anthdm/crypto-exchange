package main

import (
	"fmt"
	"time"

	"github.com/anthdm/crypto-exchange/client"
	"github.com/anthdm/crypto-exchange/server"
)

const ethPrice = 1281.0

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	go makeMarketSimple(c)
	time.Sleep(2 * time.Second)

	go marketOrderPlacer(c)

	select {}
}

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		buyOrder := client.PlaceOrderParams{
			UserID: 7,
			Bid:    true,
			Size:   1,
		}

		_, err := c.PlaceMarketOrder(&buyOrder)
		if err != nil {
			panic(err)
		}

		sellOrder := client.PlaceOrderParams{
			UserID: 7,
			Bid:    false,
			Size:   1,
		}

		_, err = c.PlaceMarketOrder(&sellOrder)
		if err != nil {
			panic(err)
		}

		<-ticker.C
	}
}

func seedMarket(c *client.Client) {
	currentPrice := ethPrice // asyns call to fetch the price
	priceOffset := 100.0

	bidOrder := client.PlaceOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  currentPrice - priceOffset,
		Size:   10,
	}

	_, err := c.PlaceLimitOrder(&bidOrder)
	if err != nil {
		panic(err)
	}

	askOrder := client.PlaceOrderParams{
		UserID: 8,
		Bid:    false,
		Price:  currentPrice + priceOffset,
		Size:   10,
	}

	_, err = c.PlaceLimitOrder(&askOrder)
	if err != nil {
		panic(err)
	}
}

func makeMarketSimple(c *client.Client) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		bestAsk, err := c.GetBestAsk()
		if err != nil {
			panic(err)
		}
		bestBid, err := c.GetBestBid()
		if err != nil {
			panic(err)
		}

		if bestAsk == 0 && bestBid == 0 {
			seedMarket(c)
			continue
		}

		fmt.Println("best Ask", bestAsk)
		fmt.Println("best bid", bestBid)

		<-ticker.C
	}
}
