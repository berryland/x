package zb

import (
	. "github.com/berryland/x"
	"testing"
	"time"
)

func TestWebSocketClient_SubscribeTicker(t *testing.T) {
	c := NewWebSocketClient()
	c.Connect()
	c.SubscribeTicker("btc_usdt", func(ticker Ticker) {
		println(ticker.Time)
		c.Disconnect()
	})

	time.Sleep(10 * time.Second)
}
