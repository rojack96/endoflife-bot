package bot

import "github.com/bwmarrin/discordgo"

// register command with slash
func applicationCommand() []*discordgo.ApplicationCommand {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Show help message",
		},
		{
			Name:        "product-list",
			Description: "Lista dei prodotti",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "page",
					Description: "Number of the page to display",
					Required:    false,
				},
			},
		},
		{
			Name:        "product-lts",
			Description: "Lista dei prodotti",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "product",
					Description: "Product to get LTS info",
					Required:    true,
				},
			},
		},
		{
			Name:        "product-info",
			Description: "Show paginated releases for a product",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "product",
					Description: "Product to show releases for",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "page",
					Description: "Page number",
					Required:    false,
				},
			},
		},
		{
			Name:        "product-releases",
			Description: "Get basic info about a specific product release",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "product",
					Description: "Product name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "release",
					Description: "Release version",
					Required:    true,
				},
			},
		},
	}
	return commands
}
