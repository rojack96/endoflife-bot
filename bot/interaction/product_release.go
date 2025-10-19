package interaction

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
	"github.com/rojack96/endoflife-bot/endoflife/dto"
	"go.uber.org/zap"
)

func (i *Interaction) ProductRelease() {
	product := ""
	release := ""
	opts := i.ic.ApplicationCommandData().Options
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
		i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide both product name and release version.",
			},
		})
		return
	}

	data := i.responseProductReleases(product, release)
	i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})

}

func (i *Interaction) responseProductReleases(product string, release string) *discordgo.InteractionResponseData {
	repo := endoflife.NewEndOfLifeRepository(i.log)
	service := endoflife.NewEndOfLifeService(repo, i.log)
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
		i.log.Error("failed to get product release info", zap.Error(err))
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

	color := 0x00FF00 // green
	if productInfo.EndOfSecuritySupport != nil {
		endOfSecuritySupport, _ := time.Parse("2006-01-02", *productInfo.EndOfSecuritySupport)
		if endOfSecuritySupport.Before(time.Now()) {
			color = 0xFF0000 // red
		}
	}

	embed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: fmt.Sprintf("%s %s", productInfo.Name, productInfo.Release),
		Color: color,
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
