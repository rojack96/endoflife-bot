package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Token string
}

func (d *DiscordBot) Run() {
	discord, err := discordgo.New("Bot " + d.Token)
	checkNilErr(err)

	// Registra l'handler per le interactions invece di messaggi
	discord.AddHandler(handleInteraction)

	// Registra i comandi slash
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Mostra l'help del bot",
		},
		{
			Name:        "product-list",
			Description: "Lista dei prodotti",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "page",
					Description: "Numero di pagina",
					Required:    false,
				},
			},
		},
	}

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
	switch i.ApplicationCommandData().Name {
	case "help":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello World😃",
			},
		})

	case "product-list":
		page := 1
		if len(i.ApplicationCommandData().Options) > 0 {
			page = int(i.ApplicationCommandData().Options[0].IntValue())
		}
		fmt.Print("qui")
		productsList := productList(page)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: productsList.Embeds,
			},
		})
	}
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}
