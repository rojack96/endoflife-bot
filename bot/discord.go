package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rojack96/endoflife-bot/endoflife"
)

type DiscordBot struct {
	Token string
}

func (d *DiscordBot) Run() {

	// create a session
	discord, err := discordgo.New("Bot " + d.Token)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(newMessage)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	case strings.Contains(message.Content, "!product-list"):
		parts := strings.Fields(message.Content)

		page := 1 // default
		if len(parts) > 1 {
			num, err := strconv.Atoi(parts[1])
			if err != nil || num < 1 {
				discord.ChannelMessageSend(message.ChannelID, "❌ error: page value not valid")
				return
			}
			page = num
		}

		productsList := productList(page)
		discord.ChannelMessageSendComplex(message.ChannelID, productsList)
	// add more cases if required
	case strings.Contains(message.Content, "!product"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
	case strings.Contains(message.Content, "!product") && strings.Contains(message.Content, "!release"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
	}

}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

// Message handlers
func productList(page int) *discordgo.MessageSend {
	var totalPages int = 0
	repo := endoflife.NewEndOfLifeRepository()
	service := endoflife.NewEndOfLifeService(repo)

	products, err := service.GetAllProducts()
	if err != nil {
		log.Fatal("Error fetching products:", err)
	}

	products, totalPages = Paginate(products, page, 25)
	fields := []*discordgo.MessageEmbedField{}
	for _, product := range products {

		field := &discordgo.MessageEmbedField{
			Name: product,
		}
		fields = append(fields, field)
	}

	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Products List [page " + strconv.Itoa(page) + " of " + strconv.Itoa(totalPages) + "]",
			Description: "List of products available:",
			Fields:      fields, // Paginate to show only first 25 products
		},
		},
	}

	return embed
}

func Paginate(items []string, page, pageSize int) ([]string, int) {
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
