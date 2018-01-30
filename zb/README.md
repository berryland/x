# Golang Client For [ZB](https://www.zb.com/)

## Usage
### HttpClient
```go
    ticker, err := NewHttpClient().GetTicker(ParsePair("btc_usdt"))
    //other codes
    //...
```

### WebSocketClient
```go
    c := NewWebSocketClient()
	c.Connect()
	c.SubscribeTicker("btc_usdt", func(ticker Ticker) {
		println(ticker.Time)
		c.Disconnect()
	})
```
