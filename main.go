package main

import (
	"time"

	"github.com/anthdm/crypto-exchange/client"
	"github.com/anthdm/crypto-exchange/mm"
	"github.com/anthdm/crypto-exchange/server"
)

func main() {
	go server.StartServer()
	time.Sleep(1 * time.Second)

	c := client.NewClient()

	cfg := mm.Config{
		UserID:         8,
		OrderSize:      10,
		MinSpread:      100,
		MakeInterval:   1 * time.Second,
		SeedOffset:     400,
		ExchangeClient: c,
	}
	maker := mm.NewMakerMaker(cfg)

	maker.Start()

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
