package bot

import (
	"fmt"
	"log"
	"net/url"
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
		{
			Name:        "product-info",
			Description: "Show paginated releases for a product",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "product",
					Description: "Product to show releases for",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "page",
					Description: "Page number",
					Required:    false,
				},
			},
		},
		{
			Name:        "product-releases",
			Description: "Get basic info about a specific product release",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "product",
					Description: "Product name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "release",
					Description: "Release version",
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

		case "product-info":
			product := ""
			page := 1
			opts := i.ApplicationCommandData().Options
			if len(opts) > 0 {
				// options may be in any order; iterate
				for _, o := range opts {
					if o.Name == "product" && o.StringValue() != "" {
						product = o.StringValue()
					}
					if o.Name == "page" {
						page = int(o.IntValue())
					}
				}
			}
			if product == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please provide a valid product name.",
					},
				})
				return
			}

			data := products(product, page)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: data,
			})
		case "product-releases":
			product := ""
			release := ""
			opts := i.ApplicationCommandData().Options
			if len(opts) > 0 {
				for _, o := range opts {
					if o.Name == "product" {
						product = o.StringValue()
					}
					if o.Name == "release" {
						release = o.StringValue()
					}
				}
			}

			if product == "" || release == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please provide both product name and release version.",
					},
				})
				return
			}

			data := productReleases(product, release)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: data,
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
			return
		}

		// Gestione bottoni paginazione releases: product_releases_prev_{escaped}_{page}
		if strings.HasPrefix(custom, "product_releases_prev_") || strings.HasPrefix(custom, "product_releases_next_") {
			parts := strings.Split(custom, "_")
			// expected: product, releases, prev|next, {escapedProduct}, {page}
			if len(parts) < 5 {
				return
			}
			escapedProduct := parts[3]
			pageStr := parts[4]
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				return
			}
			productName, err := url.QueryUnescape(escapedProduct)
			if err != nil {
				// fallback to escaped string if unescape fails
				productName = escapedProduct
			}

			newPage := page
			if parts[2] == "prev" {
				newPage = page - 1
			} else if parts[2] == "next" {
				newPage = page + 1
			}
			if newPage < 1 {
				newPage = 1
			}

			data := products(productName, newPage)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: data,
			})
			if err != nil {
				log.Printf("failed to update product releases message: %v", err)
			}
			return
		}
	}
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}
