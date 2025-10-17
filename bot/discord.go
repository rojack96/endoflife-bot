package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
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
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
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
