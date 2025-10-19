package interaction

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Interaction struct {
	session *discordgo.Session
	ic      *discordgo.InteractionCreate
	log     *zap.Logger
}

func NewInteraction(session *discordgo.Session, ic *discordgo.InteractionCreate) *Interaction {
	return &Interaction{session: session, ic: ic}
}

func (i *Interaction) SetLogger(logger *zap.Logger) {
	i.log = logger
}
