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

func extractDataApiError(value []byte) error {
	status, _ := json.GetString(value, "status")
	if status == "ok" {
		return nil
	}
	//TODO
	//code, _ := json.GetString(value, "err-code")
	msg, _ := json.GetString(value, "err-msg")
	return &ApiError{Code: GeneralError, Message: msg}
}
