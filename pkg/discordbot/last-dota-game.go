package discordbot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// SteamAccountCommand encapsulates the logic to retrieve information for a users steam account
type SteamAccountCommand struct{}

const discordBaseAPIURL string = "https://discordapp.com/api/v6/"

func (cmd SteamAccountCommand) execute(s *discordgo.Session, m *discordgo.MessageCreate) error {

	user := m.Message.Author
	userID := user.ID
	fmt.Printf("UserID: %+v", userID)

	getUserConnectionInfo(&userID)
	// TODO no SteamID?
	// TODO message channel you need to make sure Steam account is linked to Discord and give bot permission for this command to work, include Auth command

	message, error := getLastMatchDetails("")
	if error == nil {
		sendMessageToChannel(s, "channelID TODO", *message)
	}
	return nil
}

func (cmd SteamAccountCommand) keyword() string {
	return "!steamaccount"
}

func (cmd SteamAccountCommand) description() string {
	return "Returns the steam account information associated with the account"
}
