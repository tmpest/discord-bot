package discordbot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Bot is the encapsulation of the bot and commands it supports
type Bot struct {
	DiscordSession *discordgo.Session
	Commands       []Command
}

// NewBot creates a new DiscordBot
func NewBot(discordToken string, commands []Command) (*Bot, error) {
	dg, error := discordgo.New(fmt.Sprintf("Bot %s", discordToken))
	if error != nil {
		fmt.Println("Error creating Discord session,", error)
		return nil, error
	}

	bot := &Bot{
		DiscordSession: dg,
		Commands:       commands,
	}

	return bot, error
}

// Start enables the DiscordBot and it will run until killed via an interrupt signal
func (bot *Bot) Start() error {
	discordSession := bot.DiscordSession

	discordSession.AddHandler(bot.messageCreate)

	error := discordSession.Open()
	if error != nil {
		fmt.Println("Error opening the Discord Session!\n", error)
		return error
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Tmpest bot is now running.  Press CTRL-C to exit.")

	// Create a channel to wait and listen for the interupt
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	error = discordSession.Close()
	if error != nil {
		fmt.Println("Error closing the Discord Session!\n", error)
	}
	return error
}

func (bot *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	for _, command := range bot.Commands {
		if strings.HasPrefix(m.Content, command.keyword()) {
			command.execute(s, m)
			return
		}
	}

	if strings.HasPrefix(m.Content, "!help") {
		var stringBuilder strings.Builder
		stringBuilder.WriteString("The following are supported Bot commands:\n")
		for _, command := range bot.Commands {
			stringBuilder.WriteString(fmt.Sprintf("%+v\t%+v\n", command.keyword(), command.description()))
		}
		sendMessageToChannel(s, m.ChannelID, stringBuilder.String())
	}
}
