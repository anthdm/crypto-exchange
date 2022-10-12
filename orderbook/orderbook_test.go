package orderbook

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLastMarketTrades(t *testing.T) {
	ob := NewOrderbook()
	price := 10000.0

	sellOrder := NewOrder(false, 10, 0)
	ob.PlaceLimitOrder(price, sellOrder)

	marketOrder := NewOrder(true, 10, 0)
	matches := ob.PlaceMarketOrder(marketOrder)
	assert(t, len(matches), 1)
	match := matches[0]

	assert(t, len(ob.Trades), 1)
	trade := ob.Trades[0]
	assert(t, trade.Price, price)
	assert(t, trade.Bid, marketOrder.Bid)
	assert(t, trade.Size, match.SizeFilled)
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5, 0)
	buyOrderB := NewOrder(true, 8, 0)
	buyOrderC := NewOrder(true, 10, 0)

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)

	fmt.Println(l)
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 10, 0)
	sellOrderB := NewOrder(false, 5, 0)
	ob.PlaceLimitOrder(10_000, sellOrderA)
	ob.PlaceLimitOrder(9_000, sellOrderB)

	assert(t, len(ob.Orders), 2)
	assert(t, ob.Orders[sellOrderA.ID], sellOrderA)
	assert(t, ob.Orders[sellOrderB.ID], sellOrderB)
	assert(t, len(ob.asks), 2)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrder := NewOrder(false, 20, 0)
	ob.PlaceLimitOrder(10_000, sellOrder)

	buyOrder := NewOrder(true, 10, 0)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 10.0)
	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeFilled, 10.0)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, buyOrder.IsFilled(), true)
}

func TestPlaceMarketOrderMultiFill(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := NewOrder(true, 5, 0) // filled fully
	buyOrderB := NewOrder(true, 8, 0) // partially filled
	buyOrderC := NewOrder(true, 1, 0)
	buyOrderD := NewOrder(true, 1, 0)

	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(5_000, buyOrderD)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(10_000, buyOrderA)

	assert(t, ob.BidTotalVolume(), 15.00)

	sellOrder := NewOrder(false, 10, 0)
	matches := ob.PlaceMarketOrder(sellOrder)

	assert(t, ob.BidTotalVolume(), 5.00)
	assert(t, len(ob.bids), 2)
	assert(t, len(matches), 2)
}

func TestCancelOrderAsk(t *testing.T) {
	ob := NewOrderbook()
	sellOrder := NewOrder(false, 4, 0)
	price := 10_000.0
	ob.PlaceLimitOrder(price, sellOrder)

	assert(t, ob.AskTotalVolume(), 4.0)

	ob.CancelOrder(sellOrder)
	assert(t, ob.AskTotalVolume(), 0.0)

	_, ok := ob.Orders[sellOrder.ID]
	assert(t, ok, false)

	_, ok = ob.AskLimits[price]
	assert(t, ok, false)
}

func TestCancelOrderBid(t *testing.T) {
	ob := NewOrderbook()
	buyOrder := NewOrder(true, 4, 0)
	price := 10_000.0
	ob.PlaceLimitOrder(price, buyOrder)

	assert(t, ob.BidTotalVolume(), 4.0)

	ob.CancelOrder(buyOrder)
	assert(t, ob.BidTotalVolume(), 0.0)

	_, ok := ob.Orders[buyOrder.ID]
	assert(t, ok, false)

	_, ok = ob.BidLimits[price]
	assert(t, ok, false)
}
