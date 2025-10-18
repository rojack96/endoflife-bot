package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rojack96/endoflife-bot/bot"
)

func main() {
	// Try to load .env file if it exists
	godotenv.Load()

	// Get token from environment variable
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN environment variable is not set")
	}

	bot := &bot.DiscordBot{Token: token}
	bot.Run()
}
