package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/anthdm/crypto-exchange/client"
	"github.com/anthdm/crypto-exchange/server"
)

const (
	maxOrders = 3
)

var (
	tick   = 2 * time.Second
	myAsks = make(map[float64]int64)
	myBids = make(map[float64]int64)
)

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		marketSell := &client.PlaceOrderParams{
			UserID: 666,
			Bid:    false,
			Size:   1000,
		}

		orderResp, err := c.PlaceMarketOrder(marketSell)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketBuyOrder := &client.PlaceOrderParams{
			UserID: 666,
			Bid:    true,
			Size:   1000,
		}

		orderResp, err = c.PlaceMarketOrder(marketBuyOrder)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		<-ticker.C
	}
}

func makeMarketSimpel(c *client.Client) {
	ticker := time.NewTicker(tick)

	for {
		orders, err := c.GetOrders(7)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("=---------------------------------------")
		fmt.Printf("%+v\n", orders)
		fmt.Println("=---------------------------------------")

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
		if len(myBids) < 3 {
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
			myBids[bidLimit.Price] = bidOrderResp.OrderID
		}

		// place the ask
		if len(myAsks) < 3 {
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
			myAsks[askLimit.Price] = askOrderResp.OrderID
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
		Size:   1_000_000,
	}

	bid := &client.PlaceOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  9_000,
		Size:   1_000_000,
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

	// limitOrderParams := &client.PlaceOrderParams{
	// 	UserID: 8,
	// 	Bid:    false,
	// 	Price:  10_000,
	// 	Size:   5_000_000,
	// }

	// _, err := c.PlaceLimitOrder(limitOrderParams)
	// if err != nil {
	// 	panic(err)
	// }

	// otherLimitOrderParams := &client.PlaceOrderParams{
	// 	UserID: 666,
	// 	Bid:    false,
	// 	Price:  9_000,
	// 	Size:   500_000,
	// }

	// _, err = c.PlaceLimitOrder(otherLimitOrderParams)
	// if err != nil {
	// 	panic(err)
	// }

	// buyLimitOrder := &client.PlaceOrderParams{
	// 	UserID: 666,
	// 	Bid:    true,
	// 	Price:  11_000,
	// 	Size:   500_000,
	// }

	// _, err = c.PlaceLimitOrder(buyLimitOrder)
	// if err != nil {
	// 	panic(err)
	// }

	// //fmt.Println("placed limit order from the client => ", resp.OrderID)

	// marketOrderParams := &client.PlaceOrderParams{
	// 	UserID: 7,
	// 	Bid:    true,
	// 	Size:   1_000_000,
	// }

	// _, err = c.PlaceMarketOrder(marketOrderParams)
	// if err != nil {
	// 	panic(err)
	// }

	// bestBidPrice, err := c.GetBestBid()
	// if err != nil {
	// 	panic(err)
	// }
	// bestAskPrice, err := c.GetBestAsk()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("best bid price ", bestBidPrice)
	// fmt.Println("best bid price ", bestAskPrice)

	// //		fmt.Println("placed market order from the client => ", resp.OrderID)

	// time.Sleep(1 * time.Second)

	select {}
}
