package main

import "github.com/rojack96/endoflife-bot/bot"

func main() {
	bot := &bot.DiscordBot{Token: ""} // create a new instance of DiscordBot struct
	bot.Run()                         // call the run function of bot/bot.go
}
