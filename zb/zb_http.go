package zb

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	. "github.com/berryland/x"
	json "github.com/buger/jsonparser"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"net/http"
)

const (
	DataApiUrl  = "http://api.zb.com/data/v1/"
	TradeApiUrl = "https://trade.zb.com/api/"
)

var ApiCodes = map[uint16]ApiCode{
	1000: OK,
	1001: GeneralError,
	1002: InternalError,
	1003: AuthenticationFailed,
	1004: FundPasswordLocked,
	1005: IncorrectFundPassword,
	1006: AuthenticationAuditing,
	1007: EmptyChannel,
	1008: EmptyEvent,
	1009: Maintained,
	2001: InsufficientFund,
	2002: InsufficientFund,
	2003: InsufficientFund,
	2005: InsufficientFund,
	2006: InsufficientFund,
	2007: InsufficientFund,
	2008: InsufficientFund,
	2009: InsufficientFund,
	3001: OrderNotFound,
	3002: InvalidPrice,
	3003: InvalidAmount,
	3004: UserNotFound,
	3005: InvalidArgument,
	3006: InvalidIpAddress,
	3007: RequestTimeExpired,
	3008: TradeRecordNotFound,
	4001: Unavailable,
	4002: TooFrequent,
}

type ZbHttpClient struct {
	Client *HttpClient
}

func NewHttpClient() *ZbHttpClient {
	return &ZbHttpClient{Client: &HttpClient{Client: &http.Client{}}}
}

func (c *ZbHttpClient) GetSymbols() (map[string]SymbolConfig, error) {
	configs := map[string]SymbolConfig{}
	resp, err := c.Client.DoGet(DataApiUrl+"markets", Query{})
	if err != nil {
		return configs, err
	}

	bytes := resp.ReadBytes()
	err = extractDataApiError(bytes)
	if err != nil {
		return configs, err
	}

	json.ObjectEach(bytes, func(key []byte, value []byte, dataType json.ValueType, offset int) error {
		symbol, _ := json.ParseString(key)
		amountScale, _ := json.GetInt(value, "amountScale")
		priceScale, _ := json.GetInt(value, "priceScale")
		configs[symbol] = SymbolConfig{AmountScale: byte(amountScale), PriceScale: byte(priceScale)}
		return nil
	})
	return configs, nil
}

func (c *ZbHttpClient) GetTicker(symbol string) (Ticker, error) {
	q := Query{
		"market": symbol,
	}
	resp, err := c.Client.DoGet(DataApiUrl+"ticker", q)
	if err != nil {
		return Ticker{}, err
	}

	bytes := resp.ReadBytes()
	err = extractDataApiError(bytes)
	if err != nil {
		return Ticker{}, err
	}

	return marshalTicker(bytes), nil
}

func (c *ZbHttpClient) GetKlines(pair Pair, period string, since uint64, size uint16) ([]Kline, error) {
	var klines []Kline
	q := Query{
		"market": parseSymbol(pair),
		"type":   period,
		"since":  since,
		"size":   size,
	}
	resp, err := c.Client.DoGet(DataApiUrl+"kline", q)
	if err != nil {
		return klines, err
	}

	bytes := resp.ReadBytes()
	err = extractDataApiError(bytes)
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

func (c *ZbHttpClient) GetTrades(symbol string, since uint64) ([]Trade, error) {
	var trades []Trade
	q := Query{
		"market": symbol,
		"since":  since,
	}
	resp, err := c.Client.DoGet(DataApiUrl+"trades", q)
	if err != nil {
		return trades, err
	}

	bytes := resp.ReadBytes()
	err = extractDataApiError(bytes)
	if err != nil {
		return trades, err
	}

	json.ArrayEach(bytes, func(value []byte, dataType json.ValueType, offset int, err error) {
		id, _ := json.GetInt(value, "tid")
		tradeType, _ := json.GetString(value, "type")
		amountString, _ := json.GetString(value, "amount")
		priceString, _ := json.GetString(value, "price")
		time, _ := json.GetInt(value, "date")

		amount, _ := strconv.ParseFloat(amountString, 64)
		price, _ := strconv.ParseFloat(priceString, 64)

		trades = append(trades, Trade{Id: uint64(id), TradeType: ParseTradeType(tradeType), Price: price, Amount: amount, Time: uint64(time)})
	})

	return trades, nil
}

func (c *ZbHttpClient) GetDepth(symbol string, size uint8) (Depth, error) {
	q := Query{
		"market": symbol,
		"size":   size,
	}
	resp, err := c.Client.DoGet(DataApiUrl+"depth", q)
	if err != nil {
		return Depth{}, err
	}

	bytes := resp.ReadBytes()
	err = extractDataApiError(bytes)
	if err != nil {
		return Depth{}, err
	}

	time, _ := json.GetInt(bytes, "timestamp")
	asks, bids := marshalDepthEntries(bytes, "asks"), marshalDepthEntries(bytes, "bids")

	return Depth{Asks: asks, Bids: bids, Time: uint64(time)}, nil
}

func (c *ZbHttpClient) GetAccount(accessKey string, secretKey string) (Account, error) {
	q := Query{
		"accesskey": accessKey,
		"method":    "getAccountInfo",
	}.Encode()

	sign(q, secretKey)

	resp, err := c.Client.DoGet(TradeApiUrl+"getAccountInfo", q)
	if err != nil {
		return Account{}, err
	}

	bytes := resp.ReadBytes()
	err = extractTradeApiError(bytes)
	if err != nil {
		return Account{}, err
	}

	var assets []Asset
	result, _, _, _ := json.Get(bytes, "result")
	json.ArrayEach(result, func(value []byte, dataType json.ValueType, offset int, err error) {
		freezeString, _ := json.GetString(value, "freez")
		freeze, _ := strconv.ParseFloat(freezeString, 64)
		availableString, _ := json.GetString(value, "available")
		available, _ := strconv.ParseFloat(availableString, 64)
		coinCnName, _ := json.GetString(value, "cnName")
		coinEnName, _ := json.GetString(value, "enName")
		coinKey, _ := json.GetString(value, "key")
		coinUnit, _ := json.GetString(value, "unitTag")
		coinScale, _ := json.GetInt(value, "unitDecimal")
		assets = append(assets, Asset{Freeze: freeze, Available: available, Coin: Coin{CnName: coinCnName, EnName: coinEnName, Key: coinKey, Unit: coinUnit, Scale: uint8(coinScale)}})
	}, "coins")

	base, _, _, _ := json.Get(result, "base")
	username, _ := json.GetString(base, "username")
	tradePasswordEnabled, _ := json.GetBoolean(base, "trade_password_enabled")
	authGoogleEnabled, _ := json.GetBoolean(base, "auth_google_enabled")
	authMobileEnabled, _ := json.GetBoolean(base, "auth_mobile_enabled")

	return Account{Username: username, TradePasswordEnabled: tradePasswordEnabled, AuthGoogleEnabled: authGoogleEnabled, AuthMobileEnabled: authMobileEnabled, Assets: assets}, nil
}

func (c *ZbHttpClient) PlaceOrder(symbol string, price, amount float64, tradeType TradeType, accessKey, secretKey string) (uint64, error) {
	q := Query{
		"currency":  symbol,
		"price":     price,
		"amount":    amount,
		"tradeType": tradeType,
		"accesskey": accessKey,
		"method":    "order",
	}.Encode()

	sign(q, secretKey)

	resp, err := c.Client.DoGet(TradeApiUrl+"order", q)
	if err != nil {
		return 0, err
	}

	bytes := resp.ReadBytes()
	err = extractTradeApiError(bytes)
	if err != nil {
		return 0, err
	}

	idString, _ := json.GetString(bytes, "id")
	id, _ := strconv.ParseUint(idString, 10, 64)
	return id, nil
}

func (c *ZbHttpClient) CancelOrder(symbol string, id uint64, accessKey, secretKey string) error {
	q := Query{
		"currency":  symbol,
		"id":        id,
		"accesskey": accessKey,
		"method":    "cancelOrder",
	}.Encode()

	sign(q, secretKey)

	resp, err := c.Client.DoGet(TradeApiUrl+"cancelOrder", q)
	if err != nil {
		return err
	}

	bytes := resp.ReadBytes()
	err = extractTradeApiError(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (c *ZbHttpClient) GetOrder(symbol string, id uint64, accessKey, secretKey string) (Order, error) {
	q := Query{
		"currency":  symbol,
		"id":        strconv.FormatUint(id, 10),
		"accesskey": accessKey,
		"method":    "getOrder",
	}.Encode()

	sign(q, secretKey)

	resp, err := c.Client.DoGet(TradeApiUrl+"getOrder", q)
	if err != nil {
		return Order{}, err
	}

	bytes := resp.ReadBytes()
	err = extractTradeApiError(bytes)
	if err != nil {
		return Order{}, err
	}

	return parseOrder(bytes), nil
}

func (c *ZbHttpClient) GetOrders(symbol string, tradeType TradeType, page uint64, size uint16, accessKey, secretKey string) ([]Order, error) {
	u := getUrlToGetOrders(symbol, tradeType, page, size, accessKey, secretKey)
	resp, err := c.Client.DoGet(u.String(), nil)
	if err != nil {
		return []Order{}, err
	}

	bytes := resp.ReadBytes()
	err = extractTradeApiError(bytes)
	if err != nil {
		return []Order{}, err
	}

	var orders []Order
	json.ArrayEach(bytes, func(value []byte, dataType json.ValueType, offset int, err error) {
		orders = append(orders, parseOrder(value))
	})

	return orders, nil
}

func parseOrder(value []byte) Order {
	idString, _ := json.GetString(value, "id")
	id, _ := strconv.ParseUint(idString, 10, 64)
	currency, _ := json.GetString(value, "currency")
	price, _ := json.GetFloat(value, "price")
	status, _ := json.GetInt(value, "status")
	totalAmount, _ := json.GetFloat(value, "total_amount")
	tradeAmount, _ := json.GetFloat(value, "trade_amount")
	tradePrice, _ := json.GetFloat(value, "trade_price")
	tradeMoney, _ := json.GetFloat(value, "trade_money")
	tradeDate, _ := json.GetInt(value, "trade_date")
	tradeType, _ := json.GetInt(value, "type")
	return Order{Id: id, Price: price, Average: tradePrice, TotalAmount: totalAmount, TradeAmount: tradeAmount, TradeMoney: tradeMoney, Symbol: currency, Status: OrderStatus(status), TradeType: TradeType(tradeType), Time: uint64(tradeDate)}
}

func getUrlToGetOrders(symbol string, tradeType TradeType, page uint64, size uint16, accessKey, secretKey string) *url.URL {
	switch tradeType {
	case All:
		return getOrdersIgnoreTradeType(symbol, page, size, accessKey, secretKey)
	case Buy, Sell:
		return getOrdersNew(symbol, tradeType, page, size, accessKey, secretKey)
	default:
		panic("Unknown trade type: " + string(tradeType))
	}
}

func getOrdersIgnoreTradeType(symbol string, page uint64, size uint16, accessKey, secretKey string) *url.URL {
	q := Query{
		"currency":  symbol,
		"pageIndex": page,
		"pageSize":  size,
		"accesskey": accessKey,
		"method":    "getOrdersIgnoreTradeType",
	}.Encode()

	sign(q, secretKey)

	u := BuildUrl(TradeApiUrl+"getOrdersIgnoreTradeType", q)
	return u
}

func getOrdersNew(symbol string, tradeType TradeType, page uint64, size uint16, accessKey, secretKey string) *url.URL {
	q := Query{
		"currency":  symbol,
		"tradeType": tradeType,
		"pageIndex": page,
		"pageSize":  size,
		"accesskey": accessKey,
		"method":    "getOrdersNew",
	}.Encode()

	sign(q, secretKey)

	u := BuildUrl(TradeApiUrl+"getOrdersNew", q)
	return u
}

func sign(query Query, secretKey string) {
	query["sign"] = genSign(secretKey, query)
	query["reqTime"] = time.Now().Unix() * 1000
}

func genSign(secretKey string, params map[string]interface{}) string {
	h := hmac.New(md5.New, []byte(fmt.Sprintf("%x", sha1.Sum([]byte(secretKey)))))
	h.Write([]byte(getSortedQueryString(params)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func getSortedQueryString(params map[string]interface{}) string {
	keys := make([]string, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var kvs []string
	for _, k := range keys {
		kvs = append(kvs, fmt.Sprintf("%v=%v", k, params[k]))
	}

	return strings.Join(kvs, "&")
}

func extractDataApiError(value []byte) error {
	msg, err := json.GetString(value, "error")
	if err == json.KeyPathNotFoundError {
		return nil
	}
	return &ApiError{Code: GeneralError, Message: msg}
}

func extractTradeApiError(value []byte) error {
	code, err := json.GetInt(value, "code")
	if err == json.KeyPathNotFoundError {
		return nil
	}

	if c := getApiCode(code); c == OK {
		return nil
	} else {
		msg, _ := json.GetString(value, "message")
		return &ApiError{Code: c, Message: msg}
	}
}

func getApiCode(code int64) ApiCode {
	if c, ok := ApiCodes[uint16(code)]; ok {
		return c
	}

	return Unknown
}
