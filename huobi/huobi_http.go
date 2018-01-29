package zb

import (
	. "github.com/berryland/x"
	json "github.com/buger/jsonparser"
	"net/http"
	"strconv"
)

const (
	DataApiUrl  = "https://api.huobi.pro/market/"
	TradeApiUrl = "https://api.huobi.pro/v1/"
)

type HuobiHttpClient struct {
	client *http.Client
}

func NewHttpClient() *HuobiHttpClient {
	c := new(HuobiHttpClient)
	c.client = &http.Client{}
	return c
}

func (c *HuobiHttpClient) GetKlines(pair Pair, period string, since uint64, size uint16) ([]Kline, error) {
	var klines []Kline
	q := map[string]string{
		"symbol": parseSymbol(pair),
		"period":   period,
		"size":   strconv.FormatUint(uint64(size), 10),
	}
	resp, err := c.doGet(BuildUrl(DataApiUrl+"market/history/kline", q).String())
	if err != nil {
		return klines, err
	}

	bytes := resp.ReadBytes()
	err = extractDataError(bytes)
	if err != nil {
		return klines, err
	}

	json.ArrayEach(bytes, func(value []byte, dataType json.ValueType, offset int, err error) {
		time, _ := json.GetInt(value, "[0]")
		open, _ := json.GetFloat(value, "[1]")
		high, _ := json.GetFloat(value, "[2]")
		low, _ := json.GetFloat(value, "[3]")
		close, _ := json.GetFloat(value, "[4]")
		amount, _ := json.GetFloat(value, "[5]")
		klines = append(klines, Kline{Time: uint64(time), Open: open, High: high, Low: low, Close: close, Amount: amount})
	}, "data")

	return klines, nil
}
