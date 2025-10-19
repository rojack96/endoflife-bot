package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
	"go.uber.org/zap"
)

func (i *Interaction) ProductLts() {
	product := ""
	if len(i.ic.ApplicationCommandData().Options) > 0 {
		product = i.ic.ApplicationCommandData().Options[0].StringValue()
	}

	productLts := i.responseProductLts(product)
	i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: productLts.Embeds,
		},
	})
}

func (i *Interaction) responseProductLts(product string) *discordgo.MessageSend {
	repo := endoflife.NewEndOfLifeRepository(i.log)
	service := endoflife.NewEndOfLifeService(repo, i.log)

	productInfo, err := service.GetProductLts(product)
	if err != nil {
		i.log.Error("failed to get product LTS info", zap.Error(err))
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
			Fields: fields,   // Paginate to show only first 25 products
			Color:  0xFFA500, // Orange color in hexadecimal
		},
		},
	}

	return embed
}
