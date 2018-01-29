package zb

import (
	. "github.com/berryland/x"
	json "github.com/buger/jsonparser"
	"net/http"
)

const (
	DataApiUrl  = "https://api.huobi.pro/market/"
	TradeApiUrl = "https://api.huobi.pro/v1/"
)

var ApiCodes = map[string]ApiCode{
	"bad-argument": InvalidArgument,
}

type HuobiHttpClient struct {
	Client *HttpClient
}

func NewHttpClient() *HuobiHttpClient {
	return &HuobiHttpClient{Client: &HttpClient{Client: &http.Client{}}}
}

func (c *HuobiHttpClient) GetKlines(pair Pair, period string, since uint64, size uint16) ([]Kline, error) {
	var klines []Kline
	q := Query{
		"symbol": parseSymbol(pair),
		"period": period,
		"size":   size,
	}
	resp, err := c.Client.DoGet(DataApiUrl+"history/kline", q)
	if err != nil {
		return klines, err
	}

	bytes := resp.ReadBytes()
	err = extractDataApiError(bytes)
	if err != nil {
		return klines, err
	}

	time, _ := json.GetInt(bytes, "ts")
	json.ArrayEach(bytes, func(value []byte, dataType json.ValueType, offset int, err error) {
		open, _ := json.GetFloat(value, "open")
		high, _ := json.GetFloat(value, "high")
		low, _ := json.GetFloat(value, "low")
		close, _ := json.GetFloat(value, "close")
		amount, _ := json.GetFloat(value, "amount")
		klines = append(klines, Kline{Time: uint64(time), Open: open, High: high, Low: low, Close: close, Amount: amount})
	}, "data")

	return klines, nil
}

func (c *HuobiHttpClient) GetTicker(pair Pair) (Ticker, error) {
	q := Query{
		"symbol": parseSymbol(pair),
	}
	resp, err := c.Client.DoGet(DataApiUrl+"detail/merged", q)
	if err != nil {
		return Ticker{}, err
	}

	bytes := resp.ReadBytes()
	err = extractDataApiError(bytes)
	if err != nil {
		return Ticker{}, err
	}

	time, _ := json.GetInt(bytes, "ts")
	ticker, _, _, _ := json.Get(bytes, "tick")
	close, _ := json.GetFloat(ticker, "close")
	high, _ := json.GetFloat(ticker, "high")
	low, _ := json.GetFloat(ticker, "low")
	amount, _ := json.GetFloat(ticker, "amount")
	ask, _ := json.GetFloat(ticker, "ask", "[0]")
	bid, _ := json.GetFloat(ticker, "bid", "[0]")

	return Ticker{Amount: amount, High: high, Low: low, Last: close, Bid: bid, Ask: ask, Time: uint64(time)}, nil
}

func extractDataApiError(value []byte) error {
	status, _ := json.GetString(value, "status")
	if status == "ok" {
		return nil
	}

	code, _ := json.GetString(value, "err-code")
	msg, _ := json.GetString(value, "err-msg")
	return &ApiError{Code: getApiCode(code), Message: msg}
}

func getApiCode(code string) ApiCode {
	if c, ok := ApiCodes[code]; ok {
		return c
	}

	return Unknown
}
