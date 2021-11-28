package bans

import (
	"fmt"
	"strings"

	"benh.codes/mcytbot/db"
	"benh.codes/mcytbot/errors"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func GetBans() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		perms, _ := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if perms&discordgo.PermissionManageMessages != discordgo.PermissionManageMessages && i.ChannelID != "831587431756333056" {
			errors.ReturnAccessDenied(s, i)
			return
		}
		ids := strings.Split(i.ApplicationCommandData().Options[0].StringValue(), ",")
		content := fmt.Sprintf("Checking if `%s` is banned anywhere", strings.Join(ids, "`, `"))
		content = strings.TrimSuffix(content, ", ``")

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})

		if err != nil {
			errors.ReturnError(*sentry.CaptureException(err), s, i)
			return
		}
		msgs := make(map[string]string)

		for _, id := range ids {
			user, err := s.User(id)
			if err != nil {
				errors.ReturnError(*sentry.CaptureException(err), s, i)
				return
			}
			if user != nil {
				content = fmt.Sprintf("<a:vLoading:853377815630184501> Checking whether `%s#%s (%s)` is banned anywhere", user.Username, user.Discriminator, user.ID)
				msg, _ := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
					Content: content,
				})
				msgs[id] = msg.ID
			}
		}

		for _, id := range ids {
			servers := make(map[string]string)
			user, _ := s.User(id)
			if user != nil {
				fmt.Printf("Checking for bans for '%s#%s' (%s)\n", user.Username, user.Discriminator, user.ID)

				content = fmt.Sprintf("Servers that `%s#%s` (`%s`) is banned in:", user.Username, user.Discriminator, user.ID)

				for _, g := range s.State.Guilds {

					var allow bool
					if db.DB.QueryRow(`SELECT bans_enabled FROM guilds WHERE id = $1`, g.ID).Scan(&allow); allow {
						fmt.Println(allow)
						if ban, _ := s.GuildBan(g.ID, user.ID); ban != nil {
							if len(servers)%5 == 0 {
								s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msgs[id], &discordgo.WebhookEdit{
									Content: fmt.Sprintf("<a:vLoading:853377815630184501> Searching for bans for `%s#%s (%s)` **%d so far**", user.Username, user.Discriminator, user.ID, len(servers)),
								})
							}
							servers[g.ID] = g.Name
							tmp := fmt.Sprintf("\n- %s (%s) *(%s)* ", g.Name, g.ID, ban.Reason)
							if len(content)+len(tmp) > 2048 {
								s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msgs[id], &discordgo.WebhookEdit{
									Content: content,
								})
								msg, _ := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
									Content: fmt.Sprintf("<a:vLoading:853377815630184501> Searching for bans for `%s#%s (%s)` **%d so far**", user.Username, user.Discriminator, user.ID, len(servers)),
								})
								msgs[id] = msg.ID
								content = fmt.Sprintf("Bans for `%s#%s` (`%s`) ***(Continued)***", user.Username, user.Discriminator, user.ID)
							}
							content += tmp
						}
					}
				}
			}

			if len(servers) == 0 {
				content = fmt.Sprintf("`%s#%s` (`%s`) is not banned anywhere", user.Username, user.Discriminator, user.ID)
			}
			s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msgs[id], &discordgo.WebhookEdit{
				Content: content,
			})
		}
	}
}
