package bot

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
)

// Message handlers
func productList(page int) *discordgo.MessageSend {
	var totalPages int = 0
	repo := endoflife.NewEndOfLifeRepository()
	service := endoflife.NewEndOfLifeService(repo)

	products, err := service.GetAllProducts()
	if err != nil {
		log.Fatal("Error fetching products:", err)
	}

	products, totalPages = paginate(products, page, 25)
	fields := []*discordgo.MessageEmbedField{}
	for _, product := range products {

		field := &discordgo.MessageEmbedField{
			Name: product,
		}
		fields = append(fields, field)
	}

	footer := &discordgo.MessageEmbedFooter{
		Text: "Page " + strconv.Itoa(page) + " of " + strconv.Itoa(totalPages),
	}

	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Products List",
			Description: "List of products available:",
			Fields:      fields, // Paginate to show only first 25 products
			Footer:      footer,
		},
		},
	}

	return embed
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
				Value: "Product not found",
			},
		}
	} else {
		productInfoEndOfActiveSupport := "null"
		productInfoEndOfSecuritySupport := "null"
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
