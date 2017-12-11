package exchange

import (
	"github.com/toorop/go-bittrex"
)

type Bittrex struct {
	key    string
	secret string
	api    *bittrex.Bittrex
}

func (b *Bittrex) Start() {
	b.api = bittrex.New(b.key, b.secret)
}

func (b *Bittrex) GetOrderHistory() ([]bittrex.Order, error) {
	return b.api.GetOrderHistory("all")
}

func NewBittrex(k string, s string) *Bittrex {
	return &Bittrex{
		key:    k,
		secret: s,
	}
}
