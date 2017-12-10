package exchange

import (
	"github.com/tzapu/disco-bit/bot"
)

var discord *bot.Discord

func Start(k, s, t string) {

	discord = bot.NewDiscord(t)
	discord.Start()

	/*
		// Bittrex client
		bittrex := bittrex.New(k, s)

		// Get markets
		orders, err := bittrex.GetOrderHistory("all")
		utils.FatalIfError(err)
		spew.Dump(orders)
	*/

}
