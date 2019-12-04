package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

/*
 * TODO Support deleting a message that has been sent by the bot AND was reacted to with an "X" emoji
 * Ideally the bot would initially react to all messages with the "X" - see commented out reaction code
 * Currently adding a handler for messageReactions doesn't seem to provide the necessary information in the event
 * The message or id of the message that was reacted to isn't present and we need that to delete it. Could possibly
 * resolve the message by digging through other things. But saving that for later...
 */

func SendMessageToChannel(s *discordgo.Session, channelID string, body string) (message *discordgo.Message) {
	message, err := s.ChannelMessageSend(channelID, body)
	if err != nil {
		fmt.Println("Error sending message to channel: ", err)
	}
	return message
}
