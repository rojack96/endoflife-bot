package interaction

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
	"go.uber.org/zap"
)

func (i *Interaction) ProductList() {
	page := 1
	if len(i.ic.ApplicationCommandData().Options) > 0 {
		page = int(i.ic.ApplicationCommandData().Options[0].IntValue())
	}
	data := i.responseProductList(page)
	i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

func (i *Interaction) ProductListButton(custom string) {
	if strings.HasPrefix(custom, "products_prev_") || strings.HasPrefix(custom, "products_next_") {
		parts := strings.Split(custom, "_")
		if len(parts) < 3 {
			// fallback: ignore
			return
		}
		pageStr := parts[2]
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			i.log.Warn("malformed page number in button custom id", zap.String("custom_id", custom))
			// ignore malformed id
			return
		}

		newPage := page
		switch parts[1] {
		case "prev":
			newPage = page - 1
		case "next":
			newPage = page + 1
		default:
			// ignore unknown action
			return
		}

		if newPage < 1 {
			newPage = 1
		}

		// costruisci nuova pagina
		data := i.responseProductList(newPage)

		// aggiorna il messaggio originale (edit)
		err = i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: data,
		})
		if err != nil {
			i.log.Error("failed to update product list message", zap.Error(err))
		}
		return
	}
}

// Message handlers
func (i *Interaction) responseProductList(page int) *discordgo.InteractionResponseData {
	var totalPages int = 0
	repo := endoflife.NewEndOfLifeRepository(i.log)
	service := endoflife.NewEndOfLifeService(repo, i.log)

	products, err := service.GetAllProducts()
	if err != nil {
		i.log.Error("failed to get products list", zap.Error(err))
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
		Color:       0xADD8E6, // light blue
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
