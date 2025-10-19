package interaction

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
	"go.uber.org/zap"
)

func (i *Interaction) Products() {
	product := ""
	page := 1
	opts := i.ic.ApplicationCommandData().Options
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
		i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a valid product name.",
			},
		})
		return
	}

	data := i.responseProducts(product, page)
	i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

func (i *Interaction) ProductsButton(custom string) {
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
			i.log.Warn("malformed page number in button custom id", zap.String("custom_id", custom))
			return
		}
		productName, err := url.QueryUnescape(escapedProduct)
		if err != nil {
			// fallback to escaped string if unescape fails
			productName = escapedProduct
		}

		newPage := page
		switch parts[2] {
		case "prev":
			newPage = page - 1
		case "next":
			newPage = page + 1
		}
		if newPage < 1 {
			newPage = 1
		}

		data := i.responseProducts(productName, newPage)
		err = i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: data,
		})
		if err != nil {
			i.log.Error("failed to update product releases message", zap.Error(err))
		}
		return
	}
}

func (i *Interaction) responseProducts(product string, page int) *discordgo.InteractionResponseData {
	repo := endoflife.NewEndOfLifeRepository(i.log)
	service := endoflife.NewEndOfLifeService(repo, i.log)

	productInfo, err := service.GetProducts(product)
	if err != nil {
		i.log.Error("failed to get products list", zap.Error(err))
	}

	// costruisci tutti i campi (uno per ogni release)
	allFields := []*discordgo.MessageEmbedField{}
	for _, p := range productInfo {
		endActive := "--"
		endSecurity := "--"
		if p.EndOfActiveSupport != nil {
			endActive = *p.EndOfActiveSupport
		}
		if p.EndOfSecuritySupport != nil {
			endSecurity = *p.EndOfSecuritySupport
		}

		value := fmt.Sprintf("Released: %s\nEnd of Active Support: %s\nEnd of Security Support: %s\nLatest: %s (%s)\nLink: %s",
			p.Released,
			endActive,
			endSecurity,
			p.Latest.Version,
			p.Latest.Date,
			p.Latest.Link,
		)

		allFields = append(allFields, &discordgo.MessageEmbedField{
			Name:  p.Release,
			Value: value,
		})
	}

	// paginazione sui campi
	pageSize := 3
	totalItems := len(allFields)
	if totalItems == 0 {
		// messaggio di errore / vuoto
		embed := &discordgo.MessageEmbed{
			Type:  discordgo.EmbedTypeRich,
			Title: product,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Error",
					Value: "Product not exist or has no releases.",
				},
			},
		}
		return &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}}
	}

	totalPages := (totalItems + pageSize - 1) / pageSize
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalItems {
		end = totalItems
	}
	pageFields := allFields[start:end]

	footer := &discordgo.MessageEmbedFooter{
		Text: "Page " + strconv.Itoa(page) + " of " + strconv.Itoa(totalPages),
	}

	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       product,
		Description: "List of versions:",
		Fields:      pageFields,
		Footer:      footer,
	}

	// crea bottoni Prev / Next; disabilita se siamo ai limiti
	escaped := url.QueryEscape(product) // evita problemi nel custom id
	prevBtn := discordgo.Button{
		CustomID: fmt.Sprintf("product_releases_prev_%s_%d", escaped, page),
		Label:    "Prev",
		Style:    discordgo.SecondaryButton,
		Disabled: page <= 1,
	}
	nextBtn := discordgo.Button{
		CustomID: fmt.Sprintf("product_releases_next_%s_%d", escaped, page),
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
