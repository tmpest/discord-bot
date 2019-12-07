package discordbot

import (
	"github.com/bwmarrin/discordgo"
)

// Command is the interface used to encapsulate the logic for a command supported by the bot
type Command interface {
	keyword() string
	description() string
	execute(s *discordgo.Session, m *discordgo.MessageCreate) error
}
