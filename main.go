package main

import "github.com/rojack96/endoflife-bot/bot"

func main() {
	bot := &bot.DiscordBot{Token: "MTQyODc5Mjk5ODEwMzA4OTE2Mg.GiKBMw.iNwvqkCt30-c2g2yd8pNdvEdzIUD7KKwL_2bpI"} // create a new instance of DiscordBot struct
	bot.Run()                                                                                                 // call the run function of bot/bot.go
}
