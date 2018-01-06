package exchange

import (
	"fmt"
	"time"

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
		t := b.formatOrderMessage(&order)
		b.send(t)
	}

	go b.monitor()
}

func (b *Bittrex) formatOrderMessage(order *bittrex.Order) string {
	return fmt.Sprintf(`%s %s %s %s * %s = %s`, order.TimeStamp, order.OrderType, order.Exchange, order.Quantity, order.PricePerUnit, order.Price)
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
			t := fmt.Sprintf(`%s: %s %s %s for %s`, o.TimeStamp, o.Exchange, o.OrderType, o.Quantity, o.Price)
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
