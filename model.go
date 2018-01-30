package x

import (
	"strings"
)

type Currency struct {
	Symbol string
}

func ParseCurrency(currency string) Currency {
	return Currency{Symbol: currency}
}

type Pair struct {
	Base      Currency
	Valuation Currency
}

func ParsePair(pair string) Pair {
	s := strings.Split(pair, "_")
	return Pair{Base: ParseCurrency(s[0]), Valuation: ParseCurrency(s[1])}
}

type SymbolConfig struct {
	AmountScale byte
	PriceScale  byte
}

type Ticker struct {
	Amount float64
	Last   float64
	Ask    float64
	Bid    float64
	High   float64
	Low    float64
	Time   uint64
}

type Kline struct {
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Amount float64
	Time   uint64
}

type Trade struct {
	Id        uint64
	TradeType TradeType
	Price     float64
	Amount    float64
	Time      uint64
}

type TradeType int8

const (
	All  TradeType = iota - 1
	Sell
	Buy
)

func ParseTradeType(string string) TradeType {
	switch string {
	case "buy":
		return Buy
	case "sell":
		return Sell
	default:
		panic("Unknown trade type: " + string)
	}
}

type Depth struct {
	Asks []DepthEntry
	Bids []DepthEntry
	Time uint64
}

type DepthEntry struct {
	Price  float64
	Amount float64
}

type Account struct {
	Username             string
	TradePasswordEnabled bool
	AuthGoogleEnabled    bool
	AuthMobileEnabled    bool
	Assets               []Asset
}

type Asset struct {
	Freeze    float64
	Available float64
	Coin      Coin
}

type Coin struct {
	CnName string
	EnName string
	Key    string
	Unit   string
	Scale  uint8
}

type Order struct {
	Id          uint64
	Price       float64
	Average     float64
	TotalAmount float64
	TradeAmount float64
	TradeMoney  float64
	Symbol      string
	Status      OrderStatus
	TradeType   TradeType
	Time        uint64
}

type OrderStatus uint8

const (
	Pending         OrderStatus = iota
	Cancelled
	Finished
	PartiallyFilled
)
