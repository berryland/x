package zb

import (
	. "github.com/berryland/x"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	accessKey = os.Getenv("ZB_ACCESS_KEY")
	secretKey = os.Getenv("ZB_SECRET_KEY")
)

func TestZbHttpClient_GetSymbols(t *testing.T) {
	NewHttpClient().GetSymbols()
}

func TestZbHttpClient_GetTicker(t *testing.T) {
	ticker, err := NewHttpClient().GetTicker(ParsePair("btc_usdt"))
	assert.Nil(t, err)
	assert.True(t, ticker.Last > 0)
}

func TestZbHttpClient_GetKlines(t *testing.T) {
	klines, err := NewHttpClient().GetKlines(ParsePair("btc_usdt"), "5min", 1516029900000, 20)
	assert.Nil(t, err)
	assert.Equal(t, 20, len(klines))
	assert.True(t, klines[0].High > 0)
}

func TestZbHttpClient_GetTrades(t *testing.T) {
	trades, err := NewHttpClient().GetTrades("btc_usdt", 0)
	assert.Nil(t, err)
	assert.True(t, trades[0].Price > 0)
}

func TestZbHttpClient_GetDepth(t *testing.T) {
	depth, err := NewHttpClient().GetDepth("btc_usdt", 10)
	assert.Nil(t, err)
	assert.NotNil(t, depth)
	assert.True(t, depth.Time > 0)
}

func TestZbHttpClient_GetAccount(t *testing.T) {
	account, err := NewHttpClient().GetAccount(accessKey, secretKey)
	assert.Nil(t, err)
	assert.NotNil(t, account.Username)
}

func TestZbHttpClient_GetOrders(t *testing.T) {
	orders, err := NewHttpClient().GetOrders("btc_usdt", All, 0, 10, accessKey, secretKey)
	assert.Nil(t, err)
	assert.NotEmpty(t, orders)
}

func TestZbHttpClient_GetOrder(t *testing.T) {
	NewHttpClient().GetOrder("btc_usdt", 2018012160893558, accessKey, secretKey)
}

func TestZbHttpClient_PlaceOrder(t *testing.T) {
	NewHttpClient().PlaceOrder("btc_usdt", 15000, 0.01, Sell, accessKey, secretKey)
}

func TestZbHttpClient_CancelOrder(t *testing.T) {
	NewHttpClient().CancelOrder("btc_usdt", 2018012261281063, accessKey, secretKey)
}
