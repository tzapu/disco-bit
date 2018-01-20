package exchange

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	log "github.com/sirupsen/logrus"
	"github.com/tzapu/disco-bit/utils"

	"github.com/toorop/go-bittrex"
)

type Bittrex struct {
	key    string
	secret string
	id     string
	api    *bittrex.Bittrex
	orders map[string]*bittrex.Order
	last   string
	sender chan<- utils.Message
}

func (b *Bittrex) Start() {
	b.api = bittrex.New(b.key, b.secret)
	b.populateOrders()
	if b.last != "" {
		order := *b.orders[b.last]
		log.Println(order, b.last)
		h := "**Your last order**"
		b.send(h)
		p := b.getProfit(order.Exchange, order.PricePerUnit)
		t := b.formatOrderMessage(&order, p)
		b.send(t)
	}

	go b.monitor()
}

func (b *Bittrex) getProfit(pair string, price decimal.Decimal) string {
	for _, o := range b.orders {
		if o.Exchange == pair && strings.Contains(o.OrderType, "_BUY") {
			log.Debug("found order", o)
			h, _ := decimal.NewFromString("100")
			return o.PricePerUnit.Mul(h).Div(price).String()
		}
	}
	return "0"
}

func (b *Bittrex) formatOrderMessage(order *bittrex.Order, profit string) string {
	p := ""
	if profit != "0" && strings.Contains(order.OrderType, "_SELL") {
		p = fmt.Sprintf("%s%%", profit)
	}
	return fmt.Sprintf(
		`%s %s %s * %sbtc = %s %s(%s)`,
		order.Exchange,
		order.OrderType,
		order.Quantity,
		order.PricePerUnit,
		order.Price,
		p, order.TimeStamp,
	)
}

func (b *Bittrex) send(t string) {
	m := utils.Message{
		ID:   b.id,
		Text: t,
	}
	b.sender <- m
}

func (b *Bittrex) monitor() {
	for _ = range time.Tick(time.Second * 30) {
		orders, err := b.GetOrderHistory()
		if err != nil {
			log.Error(err)
			continue
		}
		for i := range orders {
			o := orders[i]
			if _, ok := b.orders[o.OrderUuid]; ok {
				continue
			}
			b.orders[o.OrderUuid] = &o
			p := b.getProfit(o.Exchange, o.PricePerUnit)
			t := b.formatOrderMessage(&o, p)
			b.send(t)
		}
		b.last = orders[0].OrderUuid
	}
}

func (b *Bittrex) populateOrders() error {
	orders, err := b.GetOrderHistory()
	if err != nil {
		return err
	}
	for i := range orders {
		o := orders[i]
		b.orders[o.OrderUuid] = &o
	}
	if len(orders) > 0 {
		b.last = orders[0].OrderUuid
	}
	return nil
}

func (b *Bittrex) GetOrderHistory() ([]bittrex.Order, error) {
	return b.api.GetOrderHistory("all")
}

func NewBittrex(key, secret, id string, sender chan<- utils.Message) *Bittrex {
	return &Bittrex{
		key:    key,
		secret: secret,
		id:     id,
		orders: map[string]*bittrex.Order{},
		last:   "",
		sender: sender,
	}
}
