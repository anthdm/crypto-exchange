package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/anthdm/crypto-exchange/client"
	"github.com/anthdm/crypto-exchange/server"
)

var (
	tick = 2 * time.Second
)

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		trades, err := c.GetTrades("ETH")
		if err != nil {
			panic(err)
		}

		if len(trades) > 0 {
			fmt.Printf("exchange price => %.2f\n", trades[len(trades)-1].Price)
		}

		otherMarketSell := &client.PlaceOrderParams{
			UserID: 8,
			Bid:    false,
			Size:   1000,
		}
		orderResp, err := c.PlaceMarketOrder(otherMarketSell)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketSell := &client.PlaceOrderParams{
			UserID: 666,
			Bid:    false,
			Size:   100,
		}
		orderResp, err = c.PlaceMarketOrder(marketSell)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketBuyOrder := &client.PlaceOrderParams{
			UserID: 666,
			Bid:    true,
			Size:   100,
		}
		orderResp, err = c.PlaceMarketOrder(marketBuyOrder)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		<-ticker.C
	}
}

const userID = 7

func makeMarketSimpel(c *client.Client) {
	ticker := time.NewTicker(tick)

	for {
		orders, err := c.GetOrders(userID)
		if err != nil {
			log.Println(err)
		}

		bestAsk, err := c.GetBestAsk()
		if err != nil {
			log.Println(err)
		}
		bestBid, err := c.GetBestBid()
		if err != nil {
			log.Println(err)
		}

		spread := math.Abs(bestBid - bestAsk)
		fmt.Println("exchange spread", spread)

		// place the bid
		if len(orders.Bids) < 3 {
			bidLimit := &client.PlaceOrderParams{
				UserID: 7,
				Bid:    true,
				Price:  bestBid + 100,
				Size:   1000,
			}

			bidOrderResp, err := c.PlaceLimitOrder(bidLimit)
			if err != nil {
				log.Println(bidOrderResp.OrderID)
			}
		}

		// place the ask
		if len(orders.Asks) < 3 {
			askLimit := &client.PlaceOrderParams{
				UserID: 7,
				Bid:    false,
				Price:  bestAsk - 100,
				Size:   1000,
			}

			askOrderResp, err := c.PlaceLimitOrder(askLimit)
			if err != nil {
				log.Println(askOrderResp.OrderID)
			}
		}

		fmt.Println("best ask price", bestAsk)
		fmt.Println("best bid price", bestBid)

		<-ticker.C
	}
}

func seedMarket(c *client.Client) error {
	ask := &client.PlaceOrderParams{
		UserID: 8,
		Bid:    false,
		Price:  10_000,
		Size:   1_0000,
	}

	bid := &client.PlaceOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  9_000,
		Size:   1_0000,
	}

	_, err := c.PlaceLimitOrder(ask)
	if err != nil {
		return err
	}

	_, err = c.PlaceLimitOrder(bid)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	if err := seedMarket(c); err != nil {
		panic(err)
	}

	go makeMarketSimpel(c)

	time.Sleep(1 * time.Second)

	marketOrderPlacer(c)

	select {}
}
