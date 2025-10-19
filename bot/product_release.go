package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
	"github.com/rojack96/endoflife-bot/endoflife/dto"
)

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
