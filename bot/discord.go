package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/bot/interaction"
)

type DiscordBot struct {
	Token string
}

func (d *DiscordBot) Run() {
	discord, err := discordgo.New("Bot " + d.Token)
	checkNilErr(err)

	// Registra l'handler per le interactions invece di messaggi
	discord.AddHandler(handleInteraction)

	commands := applicationCommand()

	discord.Open()
	defer discord.Close()

	// Registra i comandi per il bot
	for _, cmd := range commands {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", cmd)
		if err != nil {
			log.Printf("Errore creazione comando %v: %v", cmd.Name, err)
		}
	}

	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	inter := interaction.NewInteraction(s, i)
	// handle command slash interactions
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "help":
			inter.Help()
		case "product-list":
			inter.ProductList()
		case "product-lts":
			inter.ProductLts()
		case "product-info":
			inter.Products()
		case "product-releases":
			inter.ProductRelease()
		}
		return
	}

	// handle button interactions
	if i.Type == discordgo.InteractionMessageComponent {
		custom := i.MessageComponentData().CustomID
		// products_prev_{page} o products_next_{page}
		inter.ProductListButton(custom)

		// product_releases_prev_{escaped}_{page}
		inter.ProductsButton(custom)
	}
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}
