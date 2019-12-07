package main

import "github.com/tmpest/discord-bot/pkg/discordbot"

func test() {
	accountID := "1"
	discordbot.GetTokenInfo(&accountID)
}
