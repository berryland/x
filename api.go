package x

type HttpApiClient interface {
	GetKlines(pair Pair, period string, since uint64, size uint16) ([]Kline, error)
}

type WsApiClient interface {
}
