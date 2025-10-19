package interaction

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
	"github.com/rojack96/endoflife-bot/endoflife/dto"
)

func ProductRelease(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	data := responseProductReleases(product, release)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})

}

func responseProductReleases(product string, release string) *discordgo.InteractionResponseData {
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
