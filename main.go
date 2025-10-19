package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rojack96/endoflife-bot/bot"
	"github.com/rojack96/endoflife-bot/config"
)

func main() {
	log := config.GetLogger()

	log.Info("Starting EndOfLife Bot...")
	// Try to load .env file if it exists
	godotenv.Load()

	// Get token from environment variable
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN environment variable is not set")
	}

	bot := bot.NewDiscordBot(token, log)
	bot.Run()
}
