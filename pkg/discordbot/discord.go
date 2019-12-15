package discordbot

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

func sendMessageToChannel(s *discordgo.Session, channelID string, body string) (message *discordgo.Message) {
	message, err := s.ChannelMessageSend(channelID, body)
	if err != nil {
		fmt.Println("Error sending message to channel: ", err)
	}
	return message
}

func getUserConnectionInfo(userID *string) {
	tokenInfo, error := getTokenInfo(userID)
	if error != nil {
		// TODO
		// Return
	}
	if tokenInfo == nil {
		// TODO DM user with link to authorize this bot and instructions to try again
		// TODO message channel you need to make sure Steam account is linked to Discord and give bot permission for this command to work, include Auth command
	}

	// TODO check for expired token and refresh if it is
	// Should this be part of getTokenInfo? Maybe?

	// Use Token to get user's Discord Connections -> SteamID
}
