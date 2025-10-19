package interaction

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
)

func ProductList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	page := 1
	if len(i.ApplicationCommandData().Options) > 0 {
		page = int(i.ApplicationCommandData().Options[0].IntValue())
	}
	data := responseProductList(page)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

func ProductListButton(custom string, s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		data := responseProductList(newPage)

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
}

// Message handlers
func responseProductList(page int) *discordgo.InteractionResponseData {
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

func paginate(items []string, page, pageSize int) ([]string, int) {
	if pageSize <= 0 {
		return []string{}, 0
	}

	totalItems := len(items)
	totalPages := (totalItems + pageSize - 1) / pageSize // divisione arrotondata verso l'alto

	if page < 1 {
		page = 1
	}
	if page > totalPages {
		return []string{}, totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalItems {
		end = totalItems
	}

	return items[start:end], totalPages
}
