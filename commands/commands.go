package commands

import "github.com/bwmarrin/discordgo"

func GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "checkmutuals",
			Description: "Checks mutual servers for the given Id(s)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "userids",
					Description: "User Id(s)",
					Required:    true,
				},
			},
		},
		{
			Name:        "checkbanned",
			Description: "Checks if the given user is banned in any server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "userids",
					Description: "User Id(s)",
					Required:    true,
				},
			},
		},
		{
			Name:        "config",
			Description: "Sets the config for the guild",
		},
	}
}
