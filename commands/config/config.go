package config

import (
	"benh.codes/mcytbot/db"
	"benh.codes/mcytbot/errors"
	"github.com/bwmarrin/discordgo"
)

func BanConfig() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		perms, _ := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if perms&discordgo.PermissionAdministrator != discordgo.PermissionAdministrator && i.User.ID != "126429064218017802" {
			errors.ReturnAccessDenied(s, i)
			return
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please select from the following options below as to whether mods should be able to view bans on this server via `/checkbanned`.",
				Flags:   1 << 6,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Enable Bans",
								Style:    discordgo.SuccessButton,
								Disabled: false,
								CustomID: "bans_enable",
								Emoji: discordgo.ComponentEmoji{
									Name: "✅",
								},
							},
							discordgo.Button{
								Label:    "Disable Bans",
								Style:    discordgo.DangerButton,
								Disabled: false,
								CustomID: "bans_disable",
								Emoji: discordgo.ComponentEmoji{
									Name: "🚫",
								},
							},
						},
					},
				},
			},
		})
	}
}

func EnableBans() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// send follow up & set bans to true
		db.DB.Exec(`UPDATE guilds SET bans_enabled = true WHERE id = $1`, i.GuildID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "You have **enabled** bans visibility on `/checkbanned`",
			},
		})
	}
}

func DisableBans() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// send follow up & set bans to false
		db.DB.Exec(`UPDATE guilds SET bans_enabled = false WHERE id = $1`, i.GuildID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "You have **disabled** bans visibility on `/checkbanned`",
			},
		})
	}
}
