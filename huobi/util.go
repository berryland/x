package zb

import (
	. "github.com/berryland/x"
)

func parseSymbol(pair Pair) string {
	return pair.Base.Symbol + pair.Valuation.Symbol
}
