package bot

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
	"github.com/rojack96/endoflife-bot/endoflife/dto"
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

func productLts(product string) *discordgo.MessageSend {
	repo := endoflife.NewEndOfLifeRepository()
	service := endoflife.NewEndOfLifeService(repo)

	productInfo, err := service.GetProductLts(product)
	if err != nil {
		log.Fatal("Error fetching products:", err)
	}

	fields := []*discordgo.MessageEmbedField{}

	if productInfo.Name == "" {
		fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Error",
				Value: "Product don´t have LTS releases or does not exist.",
			},
		}
	} else {
		productInfoEndOfActiveSupport := "--"
		productInfoEndOfSecuritySupport := "--"
		if productInfo.EndOfActiveSupport != nil {
			productInfoEndOfActiveSupport = *productInfo.EndOfActiveSupport
		}

		if productInfo.EndOfSecuritySupport != nil {
			productInfoEndOfSecuritySupport = *productInfo.EndOfSecuritySupport
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Release",
			Value: productInfo.Release,
		}, &discordgo.MessageEmbedField{
			Name:   "Released",
			Value:  productInfo.Released,
			Inline: true,
		}, &discordgo.MessageEmbedField{
			Name:   "End of Active Support",
			Value:  productInfoEndOfActiveSupport,
			Inline: true,
		}, &discordgo.MessageEmbedField{
			Name:   "End of Security Support",
			Value:  productInfoEndOfSecuritySupport,
			Inline: true,
		}, &discordgo.MessageEmbedField{
			Name:  "Latest Version",
			Value: productInfo.Latest.Version,
		}, &discordgo.MessageEmbedField{
			Name:  "Latest Release Date",
			Value: productInfo.Latest.Date,
		}, &discordgo.MessageEmbedField{
			Name:  "Latest Release Link",
			Value: productInfo.Latest.Link,
		})
	}

	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:  discordgo.EmbedTypeRich,
			Title: productInfo.Name + " LTS Information",
			//Description: "List of products available:",
			Fields: fields, // Paginate to show only first 25 products
		},
		},
	}

	return embed
}

func products(product string, page int) *discordgo.InteractionResponseData {
	repo := endoflife.NewEndOfLifeRepository()
	service := endoflife.NewEndOfLifeService(repo)

	productInfo, err := service.GetProducts(product)
	if err != nil {
		log.Fatal("Error fetching products:", err)
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

func productReleases(product string, release string) *discordgo.InteractionResponseData {
	repo := endoflife.NewEndOfLifeRepository()
	service := endoflife.NewEndOfLifeService(repo)
	var (
		productInfo dto.Product
		err         error
	)
	if release == "latest" || release == "" {
		productInfo, err = service.GetProductReleasesLatest(product)
	} else {
		productInfo, err = service.GetProductReleases(product, release)
	}

	if err != nil {
		return &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeRich,
					Title:       "Error",
					Description: "Failed to fetch product information",
				},
			},
		}
	}

	embed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: fmt.Sprintf("%s %s", productInfo.Name, productInfo.Release),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Released",
				Value:  productInfo.Released,
				Inline: true,
			},
			{
				Name:   "Latest Version",
				Value:  productInfo.Latest.Version,
				Inline: true,
			},
			{
				Name:  "Latest Release Link",
				Value: productInfo.Latest.Link,
			},
		},
	}

	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
	}
}

func help() *discordgo.InteractionResponseData {
	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "EndOfLife Bot Help",
		Description: "Here are all available commands:",
		Color:       0x00ff00, // Green color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "/help",
				Value: "Shows this help message with all available commands.",
			},
			{
				Name: "/product-list [page]",
				Value: "Shows a paginated list of all available products.\n" +
					"Optional: `page` - The page number to display (default: 1)",
			},
			{
				Name: "/product-lts product",
				Value: "Shows LTS (Long Term Support) information for a specific product.\n" +
					"Required: `product` - The name of the product",
			},
			{
				Name: "/product-info product",
				Value: "Shows detailed information about all releases of a product in a paginated view.\n" +
					"Required: `product` - The name of the product\n",
			},
			{
				Name: "/product-releases [product] [release]",
				Value: "Shows specific information about a product release.\n" +
					"Required: `product` - The name of the product\n" +
					"Required: `release` - The version number or 'latest' for the latest release",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "EndOfLife Bot - Track software end-of-life dates",
		},
	}

	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
	}
}

/*
  Utility function
*/

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
