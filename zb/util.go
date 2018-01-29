package zb

import (
	. "github.com/berryland/x"
	json "github.com/buger/jsonparser"
	"strconv"
)

func marshalTicker(value []byte) Ticker {
	ticker, _, _, _ := json.Get(value, "ticker")
	amountString, _ := json.GetString(ticker, "vol")
	lastString, _ := json.GetString(ticker, "last")
	sellString, _ := json.GetString(ticker, "sell")
	buyString, _ := json.GetString(ticker, "buy")
	highString, _ := json.GetString(ticker, "high")
	lowString, _ := json.GetString(ticker, "low")
	timeString, _ := json.GetString(value, "date")

	amount, _ := strconv.ParseFloat(amountString, 64)
	last, _ := strconv.ParseFloat(lastString, 64)
	sell, _ := strconv.ParseFloat(sellString, 64)
	buy, _ := strconv.ParseFloat(buyString, 64)
	high, _ := strconv.ParseFloat(highString, 64)
	low, _ := strconv.ParseFloat(lowString, 64)
	time, _ := strconv.ParseUint(timeString, 10, 64)

	return Ticker{Amount: amount, Last: last, Ask: sell, Bid: buy, High: high, Low: low, Time: time}
}

func marshalDepthEntries(value []byte, keys ...string) []DepthEntry {
	var entry []DepthEntry
	json.ArrayEach(value, func(value []byte, dataType json.ValueType, offset int, err error) {
		price, _ := json.GetFloat(value, "[0]")
		amount, _ := json.GetFloat(value, "[1]")
		entry = append(entry, DepthEntry{Price: price, Amount: amount})
	}, keys...)
	return entry
}

func parseSymbol(pair Pair) string {
	return pair.Base.Symbol + "_" + pair.Valuation.Symbol
}