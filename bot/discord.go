package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

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
			Description: "Show help message",
		},
		{
			Name:        "product-list",
			Description: "Lista dei prodotti",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "page",
					Description: "Number of the page to display",
					Required:    false,
				},
			},
		},
		{
			Name:        "product-lts",
			Description: "Lista dei prodotti",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "product",
					Description: "Product to get LTS info",
					Required:    true,
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
	// Gestione comandi slash
	if i.Type == discordgo.InteractionApplicationCommand {
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
			data := productList(page)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: data,
			})

		case "product-lts":
			// Implementa la logica per il comando product-lts qui
			product := ""
			if len(i.ApplicationCommandData().Options) > 0 {
				product = i.ApplicationCommandData().Options[0].StringValue()
			}

			productLts := productLts(product)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: productLts.Embeds,
				},
			})

		}
		return
	}

	// Gestione component interactions (bottoni / select)
	if i.Type == discordgo.InteractionMessageComponent {
		custom := i.MessageComponentData().CustomID
		// aspettarsi: products_prev_{page} o products_next_{page}
		if strings.HasPrefix(custom, "products_prev_") || strings.HasPrefix(custom, "products_next_") {
			parts := strings.Split(custom, "_")
			if len(parts) < 3 {
				// fallback: ignore
				return
			}
			pageStr := parts[2]
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				// ignore malformed id
				return
			}

			newPage := page
			if parts[1] == "prev" {
				newPage = page - 1
			} else if parts[1] == "next" {
				newPage = page + 1
			}

			if newPage < 1 {
				newPage = 1
			}

			// costruisci nuova pagina
			data := productList(newPage)

			// aggiorna il messaggio originale (edit)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: data,
			})
			if err != nil {
				log.Printf("failed to update message: %v", err)
			}
		}
		return
	}
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}
