package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anthdm/crypto-exchange/orderbook"
	"github.com/anthdm/crypto-exchange/server"
)

const Endpoint = "http://localhost:3000"

type PlaceOrderParams struct {
	UserID int64
	Bid    bool
	// Price only needed for placing LIMIT orders.
	Price float64
	Size  float64
}

type Client struct {
	*http.Client
}

func NewClient() *Client {
	return &Client{
		Client: http.DefaultClient,
	}
}

func (c *Client) GetTrades(market string) ([]*orderbook.Trade, error) {
	e := fmt.Sprintf("%s/trades/%s", Endpoint, market)
	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	trades := []*orderbook.Trade{}
	if err := json.NewDecoder(resp.Body).Decode(&trades); err != nil {
		return nil, err
	}

	return trades, nil
}

func (c *Client) GetOrders(userID int64) (*server.GetOrdersResponse, error) {
	e := fmt.Sprintf("%s/order/%d", Endpoint, userID)
	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	orders := server.GetOrdersResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (c *Client) PlaceMarketOrder(p *PlaceOrderParams) (*server.PlaceOrderResponse, error) {
	params := &server.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   server.MarketOrder,
		Bid:    p.Bid,
		Size:   p.Size,
		Market: server.MarketETH,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	e := Endpoint + "/order"
	req, err := http.NewRequest(http.MethodPost, e, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	placeOrderResponse := &server.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(placeOrderResponse); err != nil {
		return nil, err
	}

	return placeOrderResponse, nil
}

func (c *Client) GetBestAsk() (*server.Order, error) {
	e := fmt.Sprintf("%s/book/ETH/ask", Endpoint)
	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	order := &server.Order{}
	if err := json.NewDecoder(resp.Body).Decode(order); err != nil {
		return nil, err
	}

	return order, err
}

func (c *Client) GetBestBid() (*server.Order, error) {
	e := fmt.Sprintf("%s/book/ETH/bid", Endpoint)
	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	order := &server.Order{}
	if err := json.NewDecoder(resp.Body).Decode(order); err != nil {
		return nil, err
	}

	return order, err
}

func (c *Client) CancelOrder(orderID int64) error {
	e := fmt.Sprintf("%s/order/%d", Endpoint, orderID)
	req, err := http.NewRequest(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) PlaceLimitOrder(p *PlaceOrderParams) (*server.PlaceOrderResponse, error) {
	if p.Size == 0.0 {
		return nil, fmt.Errorf("size cannot be 0 when placing a limit order")
	}

	params := &server.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   server.LimitOrder,
		Bid:    p.Bid,
		Size:   p.Size,
		Price:  p.Price,
		Market: server.MarketETH,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	e := Endpoint + "/order"
	req, err := http.NewRequest(http.MethodPost, e, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	placeOrderResponse := &server.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(placeOrderResponse); err != nil {
		return nil, err
	}

	return placeOrderResponse, nil
}
