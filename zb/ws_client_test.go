package zb

import (
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

	for {
		time.Sleep(5 * time.Second)
	}
}
