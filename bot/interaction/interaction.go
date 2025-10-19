package interaction

import "github.com/bwmarrin/discordgo"

type Interaction struct {
	session *discordgo.Session
	ic      *discordgo.InteractionCreate
}

func NewInteraction(session *discordgo.Session, ic *discordgo.InteractionCreate) *Interaction {
	return &Interaction{session: session, ic: ic}
}
