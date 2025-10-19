package bot

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
)

// Message handlers
func productList(page int) *discordgo.InteractionResponseData {
	var totalPages int = 0
	repo := endoflife.NewEndOfLifeRepository()
	service := endoflife.NewEndOfLifeService(repo)

	products, err := service.GetAllProducts()
	if err != nil {
		log.Fatal("Error fetching products:", err)
	}

	productsPage, totalPages := paginate(products, page, 10)
	fields := []*discordgo.MessageEmbedField{}
	for _, product := range productsPage {
		field := &discordgo.MessageEmbedField{
			Name: "- " + product,
		}
		fields = append(fields, field)
	}

	footer := &discordgo.MessageEmbedFooter{
		Text: "Page " + strconv.Itoa(page) + " of " + strconv.Itoa(totalPages),
	}

	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Products List",
		Description: "List of products available:",
		Fields:      fields,
		Footer:      footer,
	}

	// crea bottoni Prev / Next; disabilita se siamo ai limiti
	prevBtn := discordgo.Button{
		CustomID: fmt.Sprintf("products_prev_%d", page),
		Label:    "Prev",
		Style:    discordgo.SecondaryButton,
		Disabled: page <= 1,
	}
	nextBtn := discordgo.Button{
		CustomID: fmt.Sprintf("products_next_%d", page),
		Label:    "Next",
		Style:    discordgo.PrimaryButton,
		Disabled: page >= totalPages || totalPages == 0,
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			prevBtn,
			nextBtn,
		},
	}

	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{actions},
	}
}
