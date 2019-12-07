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

	// var endpointBuilder strings.Builder
	// endpointBuilder.WriteString(discordBaseAPIURL)
	// endpointBuilder.WriteString("/users/")
	// userConnectionsRequest := http.NewRequest(http.MethodGet)

	return nil
}

func (cmd SteamAccountCommand) keyword() string {
	return "!steamaccount"
}

func (cmd SteamAccountCommand) description() string {
	return "Returns the steam account information associated with the account"
}
