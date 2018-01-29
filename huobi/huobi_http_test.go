package zb

import (
	"testing"
	. "github.com/berryland/x"
	"github.com/stretchr/testify/assert"
)

func TestHuobiHttpClient_GetKlines(t *testing.T) {
	klines, err := NewHttpClient().GetKlines(ParsePair("btc_usdt"), "1min", 0, 20)
	assert.Nil(t, err)
	assert.NotEmpty(t, klines)
}
