package mutuals

import (
	"fmt"
	"strings"

	"benh.codes/mcytbot/errors"
	"benh.codes/mcytbot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func GetMutuals() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		perms, _ := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if perms&discordgo.PermissionManageMessages == discordgo.PermissionManageMessages || i.ChannelID == "926994654241644554" {
			ids := strings.Split(i.ApplicationCommandData().Options[0].StringValue(), ",")
			content := strings.TrimSuffix(fmt.Sprintf("Searching for mutual servers for `%s`", strings.Join(ids, "`, `")), ", ``")
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				errors.ReturnError(*sentry.CaptureException(err), s, i)
			}

			msgs := make(map[string]string)

			for _, id := range ids {
				user, err := s.User(id)
				if err != nil {
					errors.ReturnError(*sentry.CaptureException(err), s, i)
				}
				if user != nil {
					content = fmt.Sprintf("<a:vLoading:853377815630184501> Searching for mutual servers for `%s#%s (%s)` (Please wait, this takes a while)", user.Username, user.Discriminator, user.ID)
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
					fmt.Printf("Checking for mutual guilds for '%s#%s' (%s)\n", user.Username, user.Discriminator, user.ID)

					content = fmt.Sprintf("Mutual servers for `%s#%s` (`%s`)", user.Username, user.Discriminator, user.ID)

					for _, g := range s.State.Guilds {
						if m := utils.GetMember(s, g.ID, user.ID); m != nil {
							if len(servers)%5 == 0 {
								s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msgs[id], &discordgo.WebhookEdit{
									Content: fmt.Sprintf("<a:vLoading:853377815630184501> Searching for mutual servers for `%s#%s (%s)` **%d so far**", user.Username, user.Discriminator, user.ID, len(servers)),
								})
							}
							servers[g.ID] = g.Name
							tmp := fmt.Sprintf("\n- %s (%s)", g.Name, g.ID)
							if len(content)+len(tmp) > 2048 {
								s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msgs[id], &discordgo.WebhookEdit{
									Content: content,
								})
								msg, _ := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
									Content: fmt.Sprintf("<a:vLoading:853377815630184501> Searching for mutual servers for `%s#%s (%s)` **%d so far**", user.Username, user.Discriminator, user.ID, len(servers)),
								})
								msgs[id] = msg.ID
								content = fmt.Sprintf("Mutual servers for `%s#%s` (`%s`) ***(Continued)***", user.Username, user.Discriminator, user.ID)
							}
							content += tmp
						}
					}

					if len(servers) == 0 {
						content = fmt.Sprintf("No mutual servers found for `%s#%s` (`%s`)", user.Username, user.Discriminator, user.ID)
					}
					s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msgs[id], &discordgo.WebhookEdit{
						Content: content,
					})
				}
			}

		} else {
			errors.ReturnAccessDenied(s, i)
		}

	}
}
