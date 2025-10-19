package bot

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/bot/interaction"
	"go.uber.org/zap"
)

type DiscordBot struct {
	token string
	log   *zap.Logger
}

func NewDiscordBot(token string, logger *zap.Logger) *DiscordBot {
	return &DiscordBot{
		token: token,
		log:   logger,
	}
}

func (d *DiscordBot) Run() {
	discord, err := discordgo.New("Bot " + d.token)
	if err != nil {
		d.log.Fatal("erro to inizialize discord session", zap.Error(err))
	}

	// Registra l'handler per le interactions invece di messaggi
	discord.AddHandler(d.handleInteraction)

	commands := applicationCommand()

	discord.Open()
	defer discord.Close()

	// Registra i comandi per il bot
	for _, cmd := range commands {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", cmd)
		if err != nil {
			d.log.Error("error creating command", zap.String("command", cmd.Name), zap.Error(err))
		}
	}

	d.log.Info("End Of Life Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func (d *DiscordBot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	inter := interaction.NewInteraction(s, i)
	inter.SetLogger(d.log)
	// handle command slash interactions
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "help":
			inter.Help()
		case "product-list":
			inter.ProductList()
		case "product-lts":
			inter.ProductLts()
		case "product-info":
			inter.Products()
		case "product-releases":
			inter.ProductRelease()
		}
		return
	}

	// handle button interactions
	if i.Type == discordgo.InteractionMessageComponent {
		custom := i.MessageComponentData().CustomID
		// products_prev_{page} o products_next_{page}
		inter.ProductListButton(custom)

		// product_releases_prev_{escaped}_{page}
		inter.ProductsButton(custom)
	}
}
