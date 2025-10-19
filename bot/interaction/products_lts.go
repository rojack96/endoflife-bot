package interaction

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
)

func (i *Interaction) ProductLts() {
	product := ""
	if len(i.ic.ApplicationCommandData().Options) > 0 {
		product = i.ic.ApplicationCommandData().Options[0].StringValue()
	}

	productLts := responseProductLts(product)
	i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: productLts.Embeds,
		},
	})
}

func responseProductLts(product string) *discordgo.MessageSend {
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
