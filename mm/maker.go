package mm

import (
	"fmt"
	"time"

	"github.com/anthdm/crypto-exchange/client"
	"github.com/sirupsen/logrus"
)

type Config struct {
	UserID         int64
	OrderSize      float64
	MinSpread      float64
	SeedOffset     float64
	ExchangeClient *client.Client
	MakeInterval   time.Duration
}

type MarketMaker struct {
	userID         int64
	orderSize      float64
	minSpread      float64
	seedOffset     float64
	exchangeClient *client.Client
	makeInterval   time.Duration
}

func NewMakerMaker(cfg Config) *MarketMaker {
	return &MarketMaker{
		userID:         cfg.UserID,
		orderSize:      cfg.OrderSize,
		minSpread:      cfg.MinSpread,
		seedOffset:     cfg.SeedOffset,
		exchangeClient: cfg.ExchangeClient,
		makeInterval:   cfg.MakeInterval,
	}
}

func (mm *MarketMaker) Start() {
	logrus.WithFields(logrus.Fields{
		"id":           mm.userID,
		"orderSize":    mm.orderSize,
		"makeInterval": mm.makeInterval,
		"minSpread":    mm.minSpread,
	}).Info("starting market maker")

	go mm.makerLoop()
}

func (mm *MarketMaker) makerLoop() {
	ticker := time.NewTicker(mm.makeInterval)

	for {
		bestBid, err := mm.exchangeClient.GetBestBid()
		if err != nil {
			logrus.Error(err)
			break
		}

		bestAsk, err := mm.exchangeClient.GetBestAsk()
		if err != nil {
			logrus.Error(err)
			break
		}

		fmt.Printf("best ask => %+v\n", bestAsk)
		fmt.Printf("best bid => %+v\n", bestBid)

		if bestAsk.Price == 0 && bestBid.Price == 0 {
			if err := mm.seedMarket(); err != nil {
				logrus.Error(err)
				break
			}
		}

		<-ticker.C
	}
}

func (mm *MarketMaker) seedMarket() error {
	currentPrice := simulateFetchCurrentETHPrice()

	logrus.WithFields(logrus.Fields{
		"currentETHPrice": currentPrice,
		"seedOffset":      mm.seedOffset,
	}).Info("orderbooks empty => seeding market!")

	bidOrder := &client.PlaceOrderParams{
		UserID: mm.userID,
		Size:   mm.orderSize,
		Bid:    true,
		Price:  currentPrice - mm.seedOffset,
	}
	_, err := mm.exchangeClient.PlaceLimitOrder(bidOrder)
	if err != nil {
		return err
	}

	askOrder := &client.PlaceOrderParams{
		UserID: mm.userID,
		Size:   mm.orderSize,
		Bid:    false,
		Price:  currentPrice + mm.seedOffset,
	}
	_, err = mm.exchangeClient.PlaceLimitOrder(askOrder)

	return err
}

// this will simulate a call to an other exchange fetching
// the current ETH price so we can offset both for a bid and ask.
func simulateFetchCurrentETHPrice() float64 {
	time.Sleep(80 * time.Millisecond)

	return 1000.0
}
