# A golang client for zb.com

[![Build Status](https://travis-ci.org/berryland/zb.svg?branch=master)](https://travis-ci.org/berryland/zb)

## Set Up
```bash
dep ensure -add github.com/berryland/zb
```

## Usage
### RestClient
```go
func TestRestClient_GetTicker(t *testing.T) {
    ticker, err := NewRestClient().GetTicker("btc_usdt")
    //other codes
    //...
}
```

### WebSocketClient
```go
func TestWebSocketClient_SubscribeTicker(t *testing.T) {
    c := NewWebSocketClient()
    c.Connect()
    c.SubscribeTicker("btc_usdt", func(ticker Ticker) {
        println(ticker.Last)
    })
}
```
