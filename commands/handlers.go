package commands

import (
	"benh.codes/mcytbot/commands/bans"
	"benh.codes/mcytbot/commands/config"
	"benh.codes/mcytbot/commands/mutuals"
	"github.com/bwmarrin/discordgo"
)

func GetCommandHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"checkmutuals": mutuals.GetMutuals(),
		"checkbanned":  bans.GetBans(),
		"config":       config.BanConfig(),
	}
}

func GetComponentHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bans_enable":  config.EnableBans(),
		"bans_disable": config.DisableBans(),
	}
}
