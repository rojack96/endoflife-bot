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
	// Gestione comandi slash
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "help":
			interaction.Help(s, i)
		case "product-list":
			interaction.ProductList(s, i)
		case "product-lts":
			interaction.ProductLts(s, i)
		case "product-info":
			interaction.Products(s, i)
		case "product-releases":
			interaction.ProductRelease(s, i)
		}
		return
	}

	// Gestione component interactions (bottoni / select)
	if i.Type == discordgo.InteractionMessageComponent {
		custom := i.MessageComponentData().CustomID
		// aspettarsi: products_prev_{page} o products_next_{page}
		interaction.ProductListButton(custom, s, i)

		// Gestione bottoni paginazione releases: product_releases_prev_{escaped}_{page}
		interaction.ProductsButton(custom, s, i)
	}
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}
