package interaction

import (
	"github.com/bwmarrin/discordgo"
)

func (i *Interaction) Help() {
	data := helpResponse()
	i.session.InteractionRespond(i.ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

func helpResponse() *discordgo.InteractionResponseData {
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
