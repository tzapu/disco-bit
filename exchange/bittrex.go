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
	key      string
	secret   string
	id       string
	api      *bittrex.Bittrex
	orders   map[string]*bittrex.Order
	orderIds []string
	sender   chan<- utils.Message
}

func (b *Bittrex) Start() {
	b.api = bittrex.New(b.key, b.secret)
	b.populateOrders()
	if len(b.orderIds) > 0 {
		o := b.orders[b.orderIds[len(b.orderIds)-1]]
		log.Println(o, o)
		h := "**Your last order**"
		b.send(h)
		p := b.getProfit(o.Exchange, o.Price)
		pu := b.getProfitPerUnit(o.Exchange, o.PricePerUnit)
		t := b.formatOrderMessage(o, p, pu)
		b.send(t)
	}

	go b.monitor()
}

func (b *Bittrex) getProfit(pair string, price decimal.Decimal) string {
	for i := len(b.orderIds) - 1; i >= 0; i-- {
		o := b.orders[b.orderIds[i]]
		if o.Exchange == pair && strings.Contains(o.OrderType, "_BUY") {
			log.Debug("found order", o)
			h, _ := decimal.NewFromString("100")
			//return o.Price.Mul(h).Div(price).StringFixed(2)
			return price.Mul(h).Div(o.Price).StringFixed(2)
		}
	}
	return "0"
}

func (b *Bittrex) getProfitPerUnit(pair string, price decimal.Decimal) string {
	for i := len(b.orderIds) - 1; i >= 0; i-- {
		o := b.orders[b.orderIds[i]]
		if o.Exchange == pair && strings.Contains(o.OrderType, "_BUY") {
			log.Debug("found order", o)
			h, _ := decimal.NewFromString("100")
			//return o.Price.Mul(h).Div(price).StringFixed(2)
			return price.Mul(h).Div(o.PricePerUnit).StringFixed(2)
		}
	}
	return "0"
}

func (b *Bittrex) formatOrderMessage(order *bittrex.Order, profit string, profitPerUnit string) string {
	p := ""
	if profit != "0" && strings.Contains(order.OrderType, "_SELL") {
		p = fmt.Sprintf("\n**%s%%** **%s%%** ", profit, profitPerUnit)
	}
	return fmt.Sprintf(
		"**%s**\n%s %s * %sbtc = %s%s@ %s",
		order.Exchange,
		order.OrderType,
		order.Quantity,
		order.PricePerUnit,
		order.Price,
		p,
		order.TimeStamp.Format("15:04 MST 02 Jan"),
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
		lastId := ""
		if len(b.orderIds) > 0 {
			lastId = b.orderIds[len(b.orderIds)-1]
		}
		orders, err := b.GetOrderHistory()
		if err != nil {
			log.Error(err)
			continue
		}

		var toAdd []*bittrex.Order
		for i := range orders {
			o := orders[i]
			if o.OrderUuid == lastId {
				break
			}
			toAdd = append([]*bittrex.Order{&o}, toAdd...)
		}
		if len(toAdd) > 0 {
			for i := range toAdd {
				o := toAdd[i]
				b.orders[o.OrderUuid] = o
				b.orderIds = append(b.orderIds, o.OrderUuid)
				p := b.getProfit(o.Exchange, o.Price)
				pu := b.getProfitPerUnit(o.Exchange, o.PricePerUnit)
				t := b.formatOrderMessage(o, p, pu)
				b.send(t)
			}
		}
	}
}

func (b *Bittrex) populateOrders() error {
	orders, err := b.GetOrderHistory()
	if err != nil {
		return err
	}

	for i := len(orders) - 1; i >= 0; i-- {
		o := orders[i]
		b.orders[o.OrderUuid] = &o
		b.orderIds = append(b.orderIds, o.OrderUuid)
	}
	return nil
}

func (b *Bittrex) GetOrderHistory() ([]bittrex.Order, error) {
	return b.api.GetOrderHistory("all")
}

func NewBittrex(key, secret, id string, sender chan<- utils.Message) *Bittrex {
	return &Bittrex{
		key:      key,
		secret:   secret,
		id:       id,
		orders:   map[string]*bittrex.Order{},
		orderIds: []string{},
		sender:   sender,
	}
}
