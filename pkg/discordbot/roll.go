package discordbot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// RollCommand encapsulates the logic to perform a random roll
type RollCommand struct{}

func (cmd RollCommand) keyword() string {
	return "!roll"
}

func (cmd RollCommand) description() string {
	return "Rolls a value between 0 - 100 unless a valid integer is provided. Example: !roll or !roll 459"
}

func (cmd RollCommand) execute(s *discordgo.Session, m *discordgo.MessageCreate) error {
	max, error := parseRollArgs(s, m)
	if error != nil {
		return error
	}

	roll := rand.Intn(max)
	s.ChannelMessageDelete(m.ChannelID, m.ID) //Delete the original message
	sendMessageToChannel(s, m.ChannelID, fmt.Sprintf("%s rolled %d out of %d", m.Author.Username, roll, max))
	return error
}

// RollAllCommand encapsulates the logic to roll for multiple users in the same voice channel
type RollAllCommand struct{}

func (cmd RollAllCommand) keyword() string {
	return "!rollall"
}

func (cmd RollAllCommand) description() string {
	return "Rolls for each person in the same Voice Channel. By default a value between 0 - 100 is rolled unless a valid integer is provided. Example: !rollall or !rollall 459"
}

func (cmd RollAllCommand) execute(s *discordgo.Session, m *discordgo.MessageCreate) error {
	max, error := parseRollArgs(s, m)
	if error != nil {
		return error
	}

	users := findRollAllParticipants(s.State, m.Author.ID)
	var stringBuilder strings.Builder
	highestRoll := 0
	winner := ""
	stringBuilder.WriteString("Rolling for Everyone! Good Luck!\n")
	members := make([]*discordgo.Member, 0)
	for _, user := range users {
		member, _ := s.State.Member(s.State.Guilds[0].ID, user)
		members = append(members, member)
	}

	for _, member := range members {
		roll := rand.Intn(max)
		if highestRoll < roll {
			highestRoll = roll
			winner = member.User.Username
		}
		stringBuilder.WriteString(fmt.Sprintf("%s rolled %d out of %d\n", member.User.Username, roll, max))
	}

	if highestRoll == 0 && winner == "" {
		fmt.Println("RollAll command invoked but un able to find participants")
	} else {
		stringBuilder.WriteString(fmt.Sprintf("%s wins with a roll of %d!", winner, highestRoll))
		sendMessageToChannel(s, m.ChannelID, stringBuilder.String())
	}

	return error
}

func parseRollArgs(s *discordgo.Session, m *discordgo.MessageCreate) (max int, err error) {
	max = 100
	var args = strings.Split(m.Content, " ")

	if len(args) > 1 {
		i, err := strconv.ParseInt(args[1], 10, 0)
		if err != nil {
			fmt.Println("Error parsing roll arguments: ", err)
			sendMessageToChannel(s, m.ChannelID, "Invalid arguments provided for the command !roll\nExample: *!roll* or *!roll 1000*")
			return 0, err
		}
		max = int(i)
	}
	return max, nil
}

func findRollAllParticipants(s *discordgo.State, targeUserID string) []string {
	var targetChannelID string
	channelIDToUserMap := make(map[string][]string)

	// Only way I found to gather this information is to iterate over the guilds and within each guild iterate over the
	// various voice states. The voice states contain the userID and channelID for everyone. From this we can build a map
	// of channelID's to the list of userID's connected to the channel.
	for _, guild := range s.Guilds {
		for _, voiceState := range guild.VoiceStates {
			users, ok := channelIDToUserMap[voiceState.ChannelID]
			if !ok {
				// first entry for the map, set users to an empty slice which will eventually be appended and added to the map
				users = make([]string, 0)
			}

			if voiceState.UserID == targeUserID {
				// capture this channelID because the list of users from this channel is the result of this method
				targetChannelID = voiceState.ChannelID
			}

			users = append(users, voiceState.UserID)
			channelIDToUserMap[voiceState.ChannelID] = users
		}
	}

	if targetChannelID == "" {
		return make([]string, 0)
	}
	return channelIDToUserMap[targetChannelID]
}
