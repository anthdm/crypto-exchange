package main

import (
	"math/rand"
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
		MinSpread:      20,
		MakeInterval:   1 * time.Second,
		SeedOffset:     40,
		ExchangeClient: c,
		PriceOffset:    10,
	}
	maker := mm.NewMakerMaker(cfg)

	maker.Start()

	time.Sleep(2 * time.Second)
	go marketOrderPlacer(c)

	select {}
}

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		randint := rand.Intn(10)
		bid := true
		if randint < 7 {
			bid = false
		}

		order := client.PlaceOrderParams{
			UserID: 7,
			Bid:    bid,
			Size:   1,
		}

		_, err := c.PlaceMarketOrder(&order)
		if err != nil {
			panic(err)
		}

		<-ticker.C
	}
}
